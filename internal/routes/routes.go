package routes

import (
	"net/http"

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
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/refresh", handlers.Refresh)
		auth.POST("/logout", handlers.Logout)
	}

	// Protected routes — middleware runs first on every route in this group
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		protected.PUT("/location", handlers.UpdateLocation)
		protected.GET("/nearby", handlers.GetNearby)
	}
}
