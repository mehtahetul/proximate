package models

import (
	"time"

	"github.com/google/uuid"
)

// this struct is basically saying "Hey GORM, whenever you interact with the users table,
// map it to this Go object"
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" means never expose in JSON
	Bio       string    `json:"bio"`
	Skills    []string  `json:"skills" gorm:"type:text[]"`
	IsVisible bool      `json:"is_visible" gorm:"default:true"`
	Latitude  *float64  `json:"-" gorm:"column:latitude"`
	Longitude *float64  `json:"-" gorm:"column:longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NearbyUser is what we return to clients from the /nearby endpoint.
// Intentionally minimal — no email, no exact coordinates.
type NearbyUser struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Bio            string    `json:"bio"`
	Skills         []string  `json:"skills"`
	DistanceMetres float64   `json:"distance_metres"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
