package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mehtahetul/proximate/internal/db"
	"github.com/mehtahetul/proximate/internal/models"
	"gorm.io/gorm"
)

// ── Create Profile ───────────────────────────────────────────────────────────

type CreateProfileRequest struct {
	Name            string   `json:"name" binding:"required"`
	Headline        string   `json:"headline" binding:"required"`
	CompanyName     string   `json:"company_name" binding:"required"`
	Bio             string   `json:"bio"`
	Skills          []string `json:"skills"`
	LinkedInURL     string   `json:"linkedin_url"`
	ProfilePhotoURL string   `json:"profile_photo_url"`
}

func CreateProfile(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in token"})
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a profile already exists for this user — POST must fail if so
	var existing models.Profile
	err = db.DB.Where("user_id = ?", userID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "profile already exists"})
		return
	}
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check existing profile"})
		return
	}

	profile := models.Profile{
		UserID:          userID,
		Name:            req.Name,
		Headline:        req.Headline,
		CompanyName:     req.CompanyName,
		Bio:             req.Bio,
		Skills:          req.Skills,
		LinkedInURL:     req.LinkedInURL,
		ProfilePhotoURL: req.ProfilePhotoURL,
	}

	if err := db.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// ── Update Profile (partial) ─────────────────────────────────────────────────

// Pointer fields let us distinguish "not sent" (nil) from "sent as empty" (&"").
type UpdateProfileRequest struct {
	Name            *string   `json:"name"`
	Headline        *string   `json:"headline"`
	CompanyName     *string   `json:"company_name"`
	Bio             *string   `json:"bio"`
	Skills          *[]string `json:"skills"`
	LinkedInURL     *string   `json:"linkedin_url"`
	ProfilePhotoURL *string   `json:"profile_photo_url"`
}

func UpdateProfile(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in token"})
		return
	}

	var profile models.Profile
	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only apply fields that were actually sent in the request
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Headline != nil {
		updates["headline"] = *req.Headline
	}
	if req.CompanyName != nil {
		updates["company_name"] = *req.CompanyName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Skills != nil {
		updates["skills"] = *req.Skills
	}
	if req.LinkedInURL != nil {
		updates["linkedin_url"] = *req.LinkedInURL
	}
	if req.ProfilePhotoURL != nil {
		updates["profile_photo_url"] = *req.ProfilePhotoURL
	}
	updates["updated_at"] = gorm.Expr("NOW()")

	if err := db.DB.Model(&profile).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile"})
		return
	}

	// Re-fetch to return the fresh state
	db.DB.Where("user_id = ?", userID).First(&profile)
	c.JSON(http.StatusOK, profile)
}

// ── Get My Profile ────────────────────────────────────────────────────────────

func GetMyProfile(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in token"})
		return
	}

	var profile models.Profile
	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ── Get Profile By User ID ───────────────────────────────────────────────────

func GetProfileByID(c *gin.Context) {
	idParam := c.Param("id")
	targetID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var profile models.Profile
	if err := db.DB.Where("user_id = ?", targetID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}