package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to show payment details
type Payment struct {
	OrderId       string         `json:"orderId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	PaymentId     string         `json:"paymentId"`
	UserId        string         `json:"userId"`
	PaymentAmount float64        `json:"amount"`
	PaymentType   string         `json:"paymentType"`
	PaymentStatus string         `json:"paymentStatus"`
	RefundAmount  float64        `json:"refundAmount"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// DB model to represent the relation of which user has placed which order
type UserPayments struct {
	UserId         string         `json:"userId"`
	PaymentId      string         `json:"paymentId"`
	OrderId        string         `json:"orderId"`
	CouponRedeemed string         `json:"couponRedeemed"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
