package context

// Add to Cart Request
type AddToCartRequest struct {
	ProductId    string  `json:"productId" validate:"required" message:"ProductId "`
	ProductCount float64 `json:"productCount" validate:"required"`
}

// Remove from cart Request
type RemoveFromCart struct {
	CartId    string `json:"cartId" validate:"required"`
	ProductId string `json:"productId" validate:"required"`
}

// Remove product from cart Request
type RemoveProduct struct {
	CartId       string  `json:"cartId" validate:"required"`
	ProductId    string  `json:"productId" validate:"required"`
	ProductCount float64 `json:"productCount" validate:"required"`
}
