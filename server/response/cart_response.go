package response

import "time"

// Cart Product Details
type CartProductResponse struct {
	CartId         string    `json:"cartId"`
	ProductId      string    `json:"productId"`
	ProductCount   float64    `json:"productCount"`
	ProductPrice   float64    `json:"productPrice"`
	ProductAddedAt time.Time `json:"productAddedAt"`
}
