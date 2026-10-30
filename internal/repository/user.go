package repository

import (
	"github.com/mehtahetul/proximate/internal/models"
	"gorm.io/gorm"
)

// UpdateLocation stores a user's current position as a PostGIS point.
// Note the argument order: ST_MakePoint takes (longitude, latitude) —
// x before y — the reverse of how people usually say coordinates.
func UpdateLocation(db *gorm.DB, userID string, lat, lng float64) error {
	return db.Exec(`
		UPDATE users
		SET location = ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
		    updated_at = NOW()
		WHERE id = ?
	`, lng, lat, userID).Error
}

// FindNearby returns users within radiusMetres of the given point.
// Returns distance in metres instead of exact coordinates — privacy by design.
func FindNearby(db *gorm.DB, lat, lng float64, radiusMetres float64, excludeUserID string) ([]models.NearbyUser, error) {
	var users []models.NearbyUser

	err := db.Raw(`
		SELECT
			u.id,
			p.name,
			p.headline,
			p.company_name,
			p.bio,
			p.skills,
			p.linkedin_url,
			p.profile_photo_url,
			ST_Distance(
				u.location,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography
			) AS distance_metres
		FROM users u
		INNER JOIN profiles p ON p.user_id = u.id
		WHERE u.is_visible = true
		  AND u.id != ?
		  AND u.location IS NOT NULL
		  -- ST_DWithin on the stored, GIST-indexed column: the index can
		  -- discard every region that can't overlap the search circle.
		  AND ST_DWithin(
				u.location,
				ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
				?
			)
		-- u.id breaks ties between users at the exact same distance, giving
		-- the stable (distance, id) order that cursor pagination relies on.
		ORDER BY distance_metres ASC, u.id ASC
	`, lng, lat, excludeUserID, lng, lat, radiusMetres).Scan(&users).Error

	return users, err
}
