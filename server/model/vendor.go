package model

import (
	"time"

	"gorm.io/gorm"
)

// DB model to store vendor information
type Vendor struct {
	VendorId      string         `json:"vendorId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	VendorName    string         `json:"vendorName"`
	VendorContact string         `json:"vendorContact"`
	Is_Active     bool           `json:"isActive"`
	GstIn         string         `json:"gstIn"`
	Street        string         `json:"street" validate:"required"`
	City          string         `json:"city"`
	State         string         `json:"state"`
	PostalCode    string         `json:"postalCode"`
	Country       string         `json:"country"`
	Description   string         `json:"description"`
	Logo          []byte         `json:"logo" gorm:"bytea"`
	BannerImage   []byte         `json:"bannerImage" gorm:"bytea"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// DB model to handler store the relation of which vendor owns which inventory
type VendorInventory struct {
	VendorId  string         `json:"vendorId"`
	ProductId string         `json:"productId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
