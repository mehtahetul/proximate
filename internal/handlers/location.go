package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehtahetul/proximate/internal/cache"
	"github.com/mehtahetul/proximate/internal/db"
	"github.com/mehtahetul/proximate/internal/repository"
)

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

	// --- Cache lookup ---
	cacheKey := cache.NearbyKey(userID, lat, lng, radius)
	if cached, ok := cache.GetNearby(c.Request.Context(), db.RedisClient, cacheKey); ok {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	// --- Cache miss: query the database ---
	users, err := repository.FindNearby(db.DB, lat, lng, radius, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch nearby users"})
		return
	}

	// Serialize to JSON, store in Redis, return to client
	response := gin.H{"users": users, "count": len(users)}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not serialize response"})
		return
	}

	cache.SetNearby(c.Request.Context(), db.RedisClient, cacheKey, string(jsonBytes))
	c.Header("X-Cache", "MISS")
	c.Data(http.StatusOK, "application/json", jsonBytes)
}
