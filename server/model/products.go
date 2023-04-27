package model

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ProductId          string         `json:"productId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	ProductName        string         `json:"productName"`
	ProductDescription string         `json:"productDescription"`
	ProductPrice       float64        `json:"productPrice"`
	ProductCategory    string         `json:"productCategory"`
	ProductBrand       string         `json:"productBrand"`
	ProductImageUrl    string         `json:"productImageUrl"`
	ProductRating      float64        `json:"productRating"`
	ProductReview      int            `json:"productReview"`
	ProductInventory   int            `json:"productInventory"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
