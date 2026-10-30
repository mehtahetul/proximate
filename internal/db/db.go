package db

import (
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// databaseURL builds the Postgres connection URL from .env values.
// url.UserPassword escapes special characters in the password, so a
// password containing '@' or '/' can't break the URL.
func databaseURL(scheme string) string {
	u := url.URL{
		Scheme:   scheme,
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

func Connect() {
	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL("postgres")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}
