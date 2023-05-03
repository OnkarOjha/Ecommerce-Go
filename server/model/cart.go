package model

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	CartId     string         `json:"cartId"`
	UserId     string         `json:"userId"`
	CartCount  int            `json:"cartCount"`
	TotalPrice float64        `json:"totalPrice" fmt:"%.2f"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type CartProducts struct {
	CartId       string         `json:"cartId" gorm:"default:uuid_generate_v4();"`
	UserId       string         `json:"userId"`
	ProductId    string         `json:"productId"`
	ProductCount float64        `json:"productCount"`
	ProductPrice float64        `json:"productPrice" fmt:"%.2f"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
