package models

import (
	"time"

	"github.com/google/uuid"
)

// this struct is basically saying "Hey GORM, whenever you interact with the users table,
// map it to this Go object"
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	IsVisible bool      `json:"is_visible" gorm:"default:true"`
	Latitude  *float64  `json:"-" gorm:"column:latitude"`
	Longitude *float64  `json:"-" gorm:"column:longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NearbyUser is what we return to clients from the /nearby endpoint.
// Intentionally minimal — no email, no exact coordinates.
type NearbyUser struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Headline        string    `json:"headline"`
	CompanyName     string    `json:"company_name"`
	Bio             string    `json:"bio"`
	Skills          []string  `json:"skills"`
	LinkedInURL     string    `json:"linkedin_url,omitempty"`
	ProfilePhotoURL string    `json:"profile_photo_url,omitempty"`
	DistanceMetres  float64   `json:"distance_metres"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type Profile struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	Name            string    `json:"name" gorm:"not null"`
	Headline        string    `json:"headline" gorm:"not null"`
	CompanyName     string    `json:"company_name" gorm:"not null"`
	Bio             string    `json:"bio"`
	Skills          []string  `json:"skills" gorm:"type:text[]"`
	LinkedInURL     string    `json:"linkedin_url,omitempty" gorm:"column:linkedin_url"`
	ProfilePhotoURL string    `json:"profile_photo_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}