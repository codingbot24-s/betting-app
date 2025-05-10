package db

import (
	"github.com/codingbot24-s/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	connstr := helpers.GetEnv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(connstr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Market{})
	return db
}