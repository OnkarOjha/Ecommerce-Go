package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	CartId       string         `json:"cartId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId       string         `json:"userId"`
	ProductCount int            `json:"productCount"`
	TotalPrice   float64        `json:"totalPrice"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type CartProducts struct {
	CartId    string `json:"cartId"`
	ProductId string `json:"productId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
