package model

import "gorm.io/gorm"

type Products struct {
	gorm.Model
	ProductId          string  `json:"productId"`
	ProductName        string  `json:"productName"`
	ProductDescription string  `json:"productDescription"`
	ProductPrice       float64 `json:"productPrice"`
	ProductCategory    string  `json:"productCategory"`
	ProductBrand       string  `json:"productBrand"`
	ProductImageUrl    string  `json:"productImageUrl"`
	ProductRating      float64 `json:"productRating"`
	ProductReview      int     `json:"productReview"`
	ProductInventory   int     `json:"productInventory"`
}
