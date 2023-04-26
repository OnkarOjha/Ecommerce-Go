package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	CartId string `json:"cartId"`
	UserId string `json:"userId"`
	ProductCount int `json:"productCount"`
	TotalPrice float64 `json:"totalPrice"`
}

type CartProducts struct {
	gorm.Model
	CartId string `json:"cartId"`
	ProductId string `json:"productId"`
}

