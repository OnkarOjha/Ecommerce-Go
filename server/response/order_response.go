package response

import (
	"time"
)

type OrderCompletionResponse struct{
	OrderId string `json:"orderId"`
	UserId string `json:"userId"`
	PaymentId string `json:"paymentId"`
	PaymentAmount float64 `json:"paymentAmount"`
	PaymentDate time.Time `json:"paymentDate"`
	CartId string `json:"cartId"`
	ProductId string `json:"productId"`
	ProductName string `json:"productName"`
	ProductCategory    string  `json:"productCategory"`
	ProductBrand       string  `json:"productBrand"`
}

