package models

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID    string    
	MarketID  string    
	Side      string    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Status    string    `gorm:"not null;default:PENDING"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	
}
