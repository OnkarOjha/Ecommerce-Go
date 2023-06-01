package response

import (
	"time"
)

// Reponse to show product order completion
type OrderCompletionResponse struct {
	OrderId         string    `json:"orderId"`
	UserId          string    `json:"userId"`
	PaymentId       string    `json:"paymentId"`
	PaymentAmount   float64   `json:"paymentAmount"`
	PaymentDate     time.Time `json:"paymentDate"`
	CartId          string    `json:"cartId"`
	OrderStatus     string    `json:"orderStatus"`
	ProductId       string    `json:"productId"`
	ProductName     string    `json:"productName"`
	ProductCategory string    `json:"productCategory"`
	ProductBrand    string    `json:"productBrand"`
}

// Reponse to show cart order completion
type CartOrderCompletionResponse struct {
	OrderId       string    `json:"orderId"`
	UserId        string    `json:"userId"`
	PaymentId     string    `json:"paymentId"`
	PaymentAmount float64   `json:"paymentAmount"`
	PaymentDate   time.Time `json:"paymentDate"`
	CartId        string    `json:"cartId"`
	OrderStatus   string    `json:"orderStatus"`
}
