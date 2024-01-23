package model

import (
	"time"

	"gorm.io/gorm"
)

// Gopher is the model representing gophers.
type Gopher struct {
	ID        uint           `gorm:"primarykey" json:"id" form:"id"`
	CreatedAt time.Time      `json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" form:"deleted_at"`
	Name      string         `json:"name" form:"name"`
	Job       string         `json:"job" form:"job"`
}
