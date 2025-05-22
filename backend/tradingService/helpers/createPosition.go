package helpers

import (
	"fmt"

	"github.com/codingbot24-s/db/models"
	"gorm.io/gorm"
)


// set the status of position to win or loss default is pending
func CreatePosition(db *gorm.DB, position models.Position) (*models.Position, error) {
	tx := db.Begin()
	if err := tx.Create(&position).Error; err != nil {
		fmt.Println("error creating position")
		tx.Rollback()
	}
	//TODO: send the amount to user service so it can be subtracted from the balance
	tx.Commit()
	// return the whole position object
	return &position, nil
	
}
