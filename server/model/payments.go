package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	PaymentId     string         `json:"paymentId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId        string         `json:"userId"`
	OrderId       string         `json:"orderId"`
	PaymentAmount float64        `json:"amount"`
	PaymentType   string         `json:"paymentType"`
	PaymentStatus string         `json:"paymentStatus"`
	RefundAmount  float64        `json:"refundAmount"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type UserPayments struct {
	UserId    string         `json:"userId"`
	PaymentId string         `json:"paymentId"`
	OrderId   string         `json:"orderId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
