package db

import (
	"fmt"
	"log"

	modles "github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load .env file

	connStr := helpers.GetEnv("DATABASE_URL")
	// Connect to DB
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&modles.User{}, &modles.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	fmt.Println("Connected to DB and migrated")
}


