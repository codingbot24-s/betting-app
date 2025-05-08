package models

import "gorm.io/gorm"

type Holding struct {
    gorm.Model
    UserID   uint    `gorm:"not null"`
    Team     string  `gorm:"not null;index:idx_user_team,unique"`
    Shares   int     `gorm:"not null"`       
    AvgBuy   float64 `gorm:"not null"`       
}
