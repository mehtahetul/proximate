package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mehtahetul/proximate/internal/handlers"
	"github.com/mehtahetul/proximate/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public routes
	auth := r.Group("/auth")
	auth.Use(middleware.RateLimit("auth", 10, time.Minute))
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.Refresh)
		auth.POST("/logout", handlers.Logout)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		protected.PUT("/location",
			middleware.RateLimit("location", 30, time.Minute),
			handlers.UpdateLocation,
		)
		protected.GET("/nearby",
			middleware.RateLimit("nearby", 30, time.Minute),
			handlers.GetNearby,
		)
	}
}
