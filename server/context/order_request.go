package context

//Order Place Request struct
type OrderRequest struct {
	CartId     string `json:"cartId" validate:"required"`
	ProductId  string `json:"productId" validate:"required"`
	CardNumber string `json:"cardNumber" validate:"required"`
	ExpMonth   string `json:"expMonth" validate:"required"`
	ExpYear    string `json:"expYear" validate:"required"`
	CVC        string `json:"cvc" validate:"required"`
}

// Cancel Order struct
type CancelOrderRequest struct {
	OrderId string `json:"orderId" validate:"required"`
}

//Order all the products from cart
type CartOrderRequest struct {
	CartId     string `json:"cartId" validate:"required"`
	CardNumber string `json:"cardNumber" validate:"required"`
	ExpMonth   string `json:"expMonth" validate:"required"`
	ExpYear    string `json:"expYear" validate:"required"`
	CVC        string `json:"cvc" validate:"required"`
}

// Vendor Order status Update
type VendorOrderStatusUpdate struct {
	OrderId     string `json:"orderId" validate:"required"`
	OrderStatus string `json:"orderStatus" validate:"required"`
}
