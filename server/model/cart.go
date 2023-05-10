package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to represent which user owns which cart
type Cart struct {
	CartId     string         `json:"cartId"`
	UserId     string         `json:"userId"`
	CartCount  int            `json:"cartCount"`
	TotalPrice float64        `json:"totalPrice" fmt:"%.2f"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// DB model to represent what are the products in the given cart
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
