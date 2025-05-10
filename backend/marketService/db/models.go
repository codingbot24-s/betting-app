package db

import (
    "time"

    "gorm.io/gorm"
)

type MarketStatus string

const (
    StatusDraft    MarketStatus = "draft"
    StatusOpen     MarketStatus = "open"
    StatusClosed   MarketStatus = "closed"
    StatusResolved MarketStatus = "resolved"
)

type Market struct {
    ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Question   string         `gorm:"type:text;not null" json:"question"`
    Category   string         `gorm:"type:varchar(100)" json:"category"` // Optional
    StartTime  time.Time      `gorm:"not null" json:"start_time"`
    EndTime    time.Time      `gorm:"not null" json:"end_time"`
    Status     MarketStatus   `gorm:"type:varchar(20);not null" json:"status"`
    Outcome    *bool          `gorm:"type:boolean" json:"outcome"` // nil = unresolved
    CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
