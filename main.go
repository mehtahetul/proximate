package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mehtahetul/proximate/internal/db"
	"github.com/mehtahetul/proximate/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Schema first: bring the database up to the latest migration
	// before anything tries to query it.
	db.RunMigrations()
	db.Connect()
	db.ConnectRedis()

	r := gin.Default()

	// Trust no proxies: ClientIP() then uses the real TCP peer address and
	// ignores X-Forwarded-For, which any client can forge. On Railway
	// (Day 16) this becomes Railway's proxy range instead of nil.
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to configure trusted proxies: %v", err)
	}

	routes.RegisterRoutes(r)
	r.Run(":8080")
}
