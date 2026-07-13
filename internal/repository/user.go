package repository

import (
	"github.com/mehtahetul/proximate/internal/models"
	"gorm.io/gorm"
)

// UpdateLocation stores a user's current lat/lng
func UpdateLocation(db *gorm.DB, userID string, lat, lng float64) error {
	return db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"latitude":   lat,
			"longitude":  lng,
			"updated_at": gorm.Expr("NOW()"),
		}).Error
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
				ST_MakePoint(u.longitude, u.latitude)::geography,
				ST_MakePoint(?, ?)::geography
			) AS distance_metres
		FROM users u
		INNER JOIN profiles p ON p.user_id = u.id
		WHERE u.is_visible = true
		  AND u.id != ?
		  AND u.latitude IS NOT NULL
		  AND u.longitude IS NOT NULL
		  AND ST_DWithin(
				ST_MakePoint(u.longitude, u.latitude)::geography,
				ST_MakePoint(?, ?)::geography,
				?
			)
		ORDER BY distance_metres ASC
	`, lng, lat, excludeUserID, lng, lat, radiusMetres).Scan(&users).Error

	return users, err
}
