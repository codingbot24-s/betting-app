package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type Order struct {
	gorm.Model
	UserID 		string 	`json:"user_id" gorm:"notnull"`
	Team 		string 	`json:"team" gorm:"notnull"`
	OrderType 	string 	`json:"order_type" gorm:"notnull"`
	Quantity 	int 	`json:"quantity" gorm:"notnull"`
	Price 		float64 	`json:"price" gorm:"notnull"`	
	Total 		float64 	`json:"total" gorm:"notnull"`
}


type OutBoxEvent struct {
	// removed uuid not working with uuid
	gorm.Model
	ID 	  uuid.UUID       `gorm:"type:uuid;primaryKey"`
	EventType string 
	Payload   datatypes.JSON
	Processed bool
	CreatedAt  time.Time
}