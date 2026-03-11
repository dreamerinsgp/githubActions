package models

import (
	"time"
)

// Item represents an item in the database
type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
