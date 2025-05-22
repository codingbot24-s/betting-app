package models

import (
	"time"

	"gorm.io/gorm"
)

// models
// user will have a role 
type User struct {
	gorm.Model                 
	Username     string        `json:"username" gorm:"unique;not null" validate:"required,min=3,max=32"`
	Password     string        `json:"password" gorm:"not null" validate:"required,min=6"`
	Email        string        `json:"email" gorm:"unique;not null" validate:"required,email"`
	Balance      float64       `json:"balance" gorm:"default:0" validate:"gte=0"`
	Role         string        `json:"role" gorm:"default:'user'" validate:"required,oneof=admin user"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID"` 

}

type Transaction struct {
	gorm.Model
	UserID          uint      `json:"user_id" gorm:"not null;index"` 
	User			User      `json:"user" gorm:"foreignKey:UserID"`
	Amount          float64   `json:"amount" gorm:"not null"`
	Description     string    `json:"description" gorm:"not null"`
	TransactionDate time.Time `json:"transaction_date" gorm:"not null"`
	TransactionID   string    `json:"transaction_id" gorm:"unique;not null"`
}
