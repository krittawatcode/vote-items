package domain

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel gorm.Model definition
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
