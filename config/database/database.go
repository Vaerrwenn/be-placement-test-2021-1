package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes DB
func InitDB() {
	// Get database string information.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("No database string found.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	DB = db
}
