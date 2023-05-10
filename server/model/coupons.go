package model

import "time"

type Coupons struct {
	CouponId    string    `json:"couponId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	CouponName  string    `json:"couponName" validate:"required"`
	CouponPrice float64   `json:"couponPrice" validate:"required"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type CouponRedemptions struct {
	CouponId   string    `json:"couponId"`
	OrderId    string    `json:"orderId"`
	RedeemedAt time.Time `json:"redeemedAt"`
}
