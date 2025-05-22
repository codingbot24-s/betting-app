package db

import (
	"time"

	
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	ID             string    `json:"id"`
	Name           string    `json:"name" validate:"required"`
	StockPrice     float64   `json:"stock_price" validate:"required,gt=0"`
	TeamSymbol     string    `json:"team_symbol" validate:"required"`
	AvailableStock int       `json:"available_stock" validate:"required,gt=0"`
	TotalStock     int       `json:"total_stock" validate:"required,gt=0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
