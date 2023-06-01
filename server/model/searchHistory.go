package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to keep track of search History
type SearchHistory struct {
	gorm.Model
	SearchId        string    `json:"searchId" gorm:"default:uuid_generate_v4()"`
	ProductId       string    `json:"productId"`
	SearchTime      time.Time `json:"searchTime"`
	SearchFrequency int       `json:"searchFrequency"`
	SearchQuery     string    `json:"searchQuery"`
}
