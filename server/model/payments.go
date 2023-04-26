package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PaymentId string `json:"paymentId"`
	UserId string `json:"userId"`
	OrderId string `json:"orderId"`
	PaymentAmount float64 `json:"amount"`
	PaymentType string `json:"paymentType"`
	PaymentStatus string `json:"paymentStatus"`
	RefundAmount float64 `json:"refundAmount"`
}

type UserPayments struct{
	gorm.Model
	UserId string `json:"userId"`
	PaymentId string `json:"paymentId"`
	OrderId string `json:"orderId"`
}
