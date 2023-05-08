package context

type OrderRequest struct {
	CartId     string `json:"cartId" validate:"required"`
	ProductId  string `json:"productId" validate:"required"`
	CardNumber string `json:"cardNumber" validate:"required"`
	ExpMonth   string `json:"expMonth" validate:"required"`
	ExpYear    string `json:"expYear" validate:"required"`
	CVC        string `json:"cvc" validate:"required"`
}

type CancelOrderRequest struct {
	OrderId string `json:"orderId" validate:"required"`
}

type CartOrderRequest struct {
	CartId     string `json:"cartId" validate:"required"`
	CardNumber string `json:"cardNumber" validate:"required"`
	ExpMonth   string `json:"expMonth" validate:"required"`
	ExpYear    string `json:"expYear" validate:"required"`
	CVC        string `json:"cvc" validate:"required"`
}
