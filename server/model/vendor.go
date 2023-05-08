package model

import (
	"time"

	"gorm.io/gorm"
)

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
	Logo          string         `json:"logo"`
	BannerImage   string         `json:"bannerImage"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
