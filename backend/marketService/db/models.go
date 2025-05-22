package db

import (
	"time"

	"github.com/codingbot24-s/common"
	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Market struct {
	ID        string              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Question  string              `gorm:"type:text;not null" json:"question"`
	Category  string              `gorm:"type:varchar(100)" json:"category"` 
	StartTime time.Time           `gorm:"not null" json:"start_time"`
	EndTime   time.Time           `gorm:"not null" json:"end_time"`
	Status    common.MarketStatus `gorm:"type:varchar(20);not null" json:"status"`
	Outcome   *bool               `gorm:"type:boolean" json:"outcome"` 
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}

type OutboxEvent struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	EventType string         `gorm:"type:varchar(100);not null"`
	Payload   datatypes.JSON `gorm:"type:jsonb;not null"`
	Processed bool           `gorm:"default:false"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}


