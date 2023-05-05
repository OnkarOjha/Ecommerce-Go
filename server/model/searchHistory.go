package model

import (
	"time"

	"gorm.io/gorm"
)

type SearchHistory struct {
	gorm.Model
	SearchId        string    `json:"searchId" gorm:"default:uuid_generate_v4()"`
	ProductId       string    `json:"productId"`
	SearchTime      time.Time `json:"searchTime"`
	SearchFrequency int       `json:"searchFrequency"`
	SearchQuery     string    `json:"searchQuery"`
}
