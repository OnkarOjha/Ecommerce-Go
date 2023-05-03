package request

type OrderRequest struct{
	CartId string `json:"cartId" validate:"required"`
	ProductId string `json:"productId" validate:"required"`
	CardNumber string `json:"cardNumber" validate:"required"`
	ExpMonth string `json:"expMonth" validate:"required"`
	ExpYear string `json:"expYear" validate:"required"`
	CVC string `json:"cvc" validate:"required"`
}