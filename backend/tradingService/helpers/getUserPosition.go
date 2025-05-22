package helpers

import (
	"time"

	"gorm.io/gorm"
)


type Position struct {
	ID        string
	MarketID  string
	Side      string
	Amount    float64
	Status    string
	CreatedAt time.Time
	
}

func GetUserPositions(db *gorm.DB, userID string) ([]Position, error) {
	var positions []Position
	result := db.Where("user_id = ?", userID).Find(&positions)
	if result.Error != nil {
		return nil, result.Error
	}
	return positions, nil
} 