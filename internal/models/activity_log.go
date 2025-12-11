package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityType string

const (
	ActivityTypeStockIncrement ActivityType = "STOCK_INCREMENT"
	ActivityTypeStockDecrement ActivityType = "STOCK_DECREMENT"
	ActivityTypeItemCreated    ActivityType = "ITEM_CREATED"
	ActivityTypeItemUpdated    ActivityType = "ITEM_UPDATED"
	ActivityTypeItemDeleted    ActivityType = "ITEM_DELETED"
)

type ActivityLog struct {
	ID          string       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string       `gorm:"not null" json:"user_id"`
	UserName    string       `gorm:"not null" json:"user_name"`
	ItemID      string       `gorm:"not null" json:"item_id"`
	ItemName    string       `gorm:"not null" json:"item_name"`
	Action      ActivityType `gorm:"not null" json:"action"`
	Quantity    int          `json:"quantity"`
	OldStock    int          `json:"old_stock"`
	NewStock    int          `json:"new_stock"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}