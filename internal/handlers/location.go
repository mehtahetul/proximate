package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehtahetul/proximate/internal/cache"
	"github.com/mehtahetul/proximate/internal/db"
	"github.com/mehtahetul/proximate/internal/models"
	"github.com/mehtahetul/proximate/internal/pagination"
	"github.com/mehtahetul/proximate/internal/repository"
	"github.com/google/uuid"
)

// defaultPageSize and maxPageSize bound how many results go out per page.
// A cap exists so a client can't request an absurd limit and force a huge
// response payload in one go.
const (
	defaultPageSize = 20
	maxPageSize     = 50
)

// NearbyResponse is the shape returned by GET /nearby. Kept separate from
// models.NearbyUser (which represents a single row) so the "page of
// results" concept doesn't leak into the per-user model.
type NearbyResponse struct {
	Users      []models.NearbyUser `json:"users"`
	Count      int                 `json:"count"`
	NextCursor string              `json:"next_cursor,omitempty"`
	HasMore    bool                `json:"has_more"`
}

// NearbyUserResponse is the client-facing shape for a single nearby user —
// same fields as models.NearbyUser but with distance bucketed instead of
// exposing the exact metre value.
type NearbyUserResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Headline        string    `json:"headline"`
	CompanyName     string    `json:"company_name"`
	Bio             string    `json:"bio"`
	Skills          []string  `json:"skills"`
	LinkedInURL     string    `json:"linkedin_url,omitempty"`
	ProfilePhotoURL string    `json:"profile_photo_url,omitempty"`
	DistanceBucket  string    `json:"distance_metres"`
}

// UpdateLocation — PUT /location
type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

func UpdateLocation(c *gin.Context) {
	// Read user ID that middleware stored on the context
	userID := c.MustGet("user_id").(string)

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.UpdateLocation(db.DB, userID, req.Latitude, req.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}

// GetNearby — GET /nearby?radius=500
// GetNearby — GET /nearby?radius=500&lat=12.34&lng=56.78
func GetNearby(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	radiusStr := c.DefaultQuery("radius", "500")
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
		return
	}
	if radius > 50000 {
		radius = 50000
	}

	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lng"})
		return
	}

	// --- Pagination params ---
	limit := defaultPageSize
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = parsedLimit
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	// A cursor marks "the last item the client already saw." If absent,
	// this is the first page — there is nothing to skip.
	var cursor *pagination.NearbyCursor
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		decoded, err := pagination.DecodeCursor(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor"})
			return
		}
		cursor = &decoded
	}

	// --- Fetch the FULL nearby list (cache-hit or cache-miss) ---
	// We deliberately cache and fetch the entire unpaginated result set,
	// then slice out the requested page in application code below. This
	// keeps every page within a single 30s cache window reading from the
	// same consistent snapshot — see repository/user.go comment for why
	// per-page cache keys would reintroduce the skip/duplicate problem
	// cursors exist to solve.
	cacheKey := cache.NearbyKey(userID, lat, lng, radius)
	var allUsers []models.NearbyUser
	cacheHit := false

	// found tracks whether Redis actually had this key — separate from
	// whether the decoded list happens to be empty. A user with zero
	// nearby results is a perfectly valid, cacheable outcome; without
	// this flag it would look identical to "cache miss" and we'd hit the
	// DB on every request for that case, silently defeating the cache.
	if cached, found := cache.GetNearby(c.Request.Context(), db.RedisClient, cacheKey); found {
		if err := json.Unmarshal([]byte(cached), &allUsers); err == nil {
			cacheHit = true
			c.Header("X-Cache", "HIT")
		}
		// If unmarshal fails, cacheHit stays false and we fall through
		// to the DB below — treating a corrupt cache entry like a miss.
	}

	if !cacheHit {
		fetched, err := repository.FindNearby(db.DB, lat, lng, radius, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch nearby users"})
			return
		}
		allUsers = fetched

		jsonBytes, err := json.Marshal(allUsers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not serialize response"})
			return
		}
		cache.SetNearby(c.Request.Context(), db.RedisClient, cacheKey, string(jsonBytes))
		c.Header("X-Cache", "MISS")
	}
	for i := range allUsers {
		allUsers[i].DistanceBucket = models.BucketDistance(allUsers[i].DistanceMetres)
	}

	// --- Paginate in Go: find where the cursor leaves off, then slice ---
	page, nextCursor, hasMore := paginateNearby(allUsers, cursor, limit)

	c.JSON(http.StatusOK, NearbyResponse{
		Users:      page,
		Count:      len(page),
		NextCursor: nextCursor,
		HasMore:    hasMore,
	})
}

// paginateNearby slices out one page from an already-sorted (by distance
// ascending, then user ID as tiebreak) list of nearby users. Returns the
// page, the cursor to fetch the next page (empty string if there isn't
// one), and whether more results exist beyond this page.
func paginateNearby(all []models.NearbyUser, cursor *pagination.NearbyCursor, limit int) ([]models.NearbyUser, string, bool) {
	startIdx := 0

	if cursor != nil {
		// Find the first row strictly after the cursor's (distance, id).
		// Linear scan is fine here — this list is already fully in memory
		// (came straight from cache or one DB query), and realistic list
		// sizes for a proximity feed don't call for anything fancier.
		for i, u := range all {
			after := u.DistanceMetres > cursor.DistanceMetres ||
				(u.DistanceMetres == cursor.DistanceMetres && u.ID.String() > cursor.UserID.String())
			if after {
				startIdx = i
				break
			}
			// If we scan the whole list without finding anything "after"
			// the cursor, startIdx should land past the end so the slice
			// below correctly comes back empty.
			startIdx = i + 1
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(all) {
		endIdx = len(all)
	}
	if startIdx > len(all) {
		startIdx = len(all)
	}

	page := all[startIdx:endIdx]

	hasMore := endIdx < len(all)
	nextCursor := ""
	if hasMore {
		last := page[len(page)-1]
		encoded, err := pagination.EncodeCursor(pagination.NearbyCursor{
			DistanceMetres: last.DistanceMetres,
			UserID:         last.ID,
		})
		if err == nil {
			nextCursor = encoded
		}
	}

	return page, nextCursor, hasMore
}