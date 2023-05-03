package model

import (
	"time"

	"gorm.io/gorm"
)

type UserAddresses struct {
	UserId      string         `json:"userId"`
	Address     string         `json:"address"`
	AddressType string         `json:"addressType"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
