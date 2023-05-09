package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to manage DB versioning
type DbVersion struct {
	Version   int            `json:"version"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
