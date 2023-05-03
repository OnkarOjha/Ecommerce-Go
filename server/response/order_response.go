package response

import (
	"github.com/gin-gonic/gin"
	"time"
)

type OrderCompletionResponse struct{
	Status  string      `json:"status"`
	Code    int64       `json:"code"`
	Message string      `json:"message"`
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

func OrderResponse(status string , statusCode int64 , message string ,orderId string , userId string , paymentId string , paymentAmount float64 , paymentDate time.Time, cartId string , productId string , productName string , productCategory string , productBrand string , context *gin.Context){
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Writer.WriteHeader(int(statusCode))

	Response(context , int(statusCode), OrderCompletionResponse{
		Status:  status,
        Code:    statusCode,
        Message: message,
        OrderId: orderId,
        UserId: userId,
        PaymentId: paymentId,
        PaymentAmount: paymentAmount,
        PaymentDate: paymentDate,
        CartId: cartId,
        ProductId: productId,
        ProductName: productName,
        ProductCategory: productCategory,
        ProductBrand: productBrand,
	})
}