package db

import (
	"log"

	"github.com/codingbot24-s/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() *gorm.DB {
	connstr := helpers.GetEnv("TEAM_SERVICE_CONECTION_STRING")

	var err error
	DB, err = gorm.Open(postgres.Open(connstr), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB.AutoMigrate(&Team{})
	return DB
}


