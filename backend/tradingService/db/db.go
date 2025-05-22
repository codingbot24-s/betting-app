package db

import (
	"fmt"
	"log"

	"github.com/codingbot24-s/db/models"
	"github.com/codingbot24-s/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectToDB() *gorm.DB {
	connStr := helpers.GetEnv("TRADING_SERVICE_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
	AutoMigrate(db)
	fmt.Println("Connected to DB and migrated")
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Order{}, &models.Holding{}, &models.OutBoxEvent{}, &models.Position{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}
}
