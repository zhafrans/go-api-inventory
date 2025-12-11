package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	MinStock    int       `gorm:"default:10" json:"min_stock"`
	MaxStock    int       `gorm:"default:100" json:"max_stock"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
	SKU         string    `gorm:"uniqueIndex" json:"sku"`
	Location    string    `json:"location"`
	CreatedBy   string    `gorm:"not null" json:"created_by"`
	Creator     *User     `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func randomString(n int) string {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        panic(err)
    }
    return hex.EncodeToString(b)
}

func (i *Item) BeforeCreate(tx *gorm.DB) error {
    i.ID = uuid.New().String()

    if i.SKU == "" {
        i.SKU = fmt.Sprintf("ITM-%d-%s",
            time.Now().Unix(),
            randomString(4),
        )
    }
    return nil
}

type CreateItemRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock" validate:"min=0"`
	MinStock    int     `json:"min_stock" validate:"min=0"`
	MaxStock    int     `json:"max_stock" validate:"min=0"`
	Price       float64 `json:"price" validate:"min=0"`
	SKU         string  `json:"sku"`
	Location    string  `json:"location"`
}

type UpdateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	MinStock    int     `json:"min_stock" validate:"min=0"`
	MaxStock    int     `json:"max_stock" validate:"min=0"`
	Price       float64 `json:"price" validate:"min=0"`
	Location    string  `json:"location"`
}

type UpdateStockRequest struct {
	Quantity int    `json:"quantity" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=increment decrement"`
	Reason   string `json:"reason"`
}