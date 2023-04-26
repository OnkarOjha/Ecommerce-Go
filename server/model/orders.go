package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderId string `json:"orderId"`
	ProductId string `json:"productId"`
	PaymentId string `json:"paymentId"`
	UserId string `json:"userId"`
	OrderStatus string `json:"orderStatus"`
	OrderAddress string `json:"orderAddress"`
	OrderDate string `json:"orderDate"`
	OrderQuantity string `json:"orderQuantity"`
}
