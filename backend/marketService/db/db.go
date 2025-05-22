package db

import (
	"fmt"
	"log"
	"time"

	"github.com/codingbot24-s/common"
	"github.com/codingbot24-s/helpers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load .env file

	connStr := helpers.GetEnv("DATABASE_URL")
	// Connect to DB
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&Market{}, &OutboxEvent{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	fmt.Println("Connected to DB and migrated")
}

func StartMarketStatusUpdater() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			var markets []Market
			if err := DB.Where("status = ?", common.StatusOpen).Find(&markets).Error; err != nil {
				log.Printf("Error fetching open markets: %v", err)
				continue
			}

			for i := range markets {
				if err := updateMarketStatus(&markets[i]); err != nil {
					log.Printf("Error updating market %s: %v", markets[i].ID, err)
				}
			}
		}
	}()
}

func updateMarketStatus(market *Market) error {
	currentTime := time.Now()
	if currentTime.After(market.EndTime) {
		market.Status = common.StatusClosed
		// print which market is closed
		fmt.Println("closing market with id ", market.ID)
		return DB.Save(market).Error
	}
	return nil
}
