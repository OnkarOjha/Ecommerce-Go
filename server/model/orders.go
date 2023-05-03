package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderId       string         `json:"orderId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	CartId        string         `json:"cartId"`
	ProductId     string         `json:"productId"`
	PaymentId     string         `json:"paymentId"`
	UserId        string         `json:"userId"`
	OrderStatus   string         `json:"orderStatus"`
	OrderAddress  string         `json:"orderAddress"`
	OrderDate     string         `json:"orderDate"`
	OrderQuantity string         `json:"orderQuantity"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
