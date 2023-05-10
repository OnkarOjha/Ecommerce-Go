package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// DB model for vendors to list their products
type Products struct {
	ProductId          string  `json:"productId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	ProductName        string  `json:"productName"`
	ProductDescription string  `json:"productDescription"`
	ProductPrice       float64 `json:"productPrice"`
	ProductCategory    string  `json:"productCategory"`
	ProductBrand       string  `json:"productBrand"`
	ProductRating      float64 `json:"productRating"`
	ProductReview      string  `json:"productReview"`
	ProductInventory   int     `json:"productInventory"`
}

func (p Products) ValidateProduct() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.ProductName, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.ProductDescription, validation.Required, validation.Length(1, 1000)),
		validation.Field(&p.ProductPrice, validation.Required),
		validation.Field(&p.ProductCategory, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.ProductBrand, validation.Required, validation.Length(1, 50)),
	)
}
