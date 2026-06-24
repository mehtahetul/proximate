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

	db.Connect()

	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
