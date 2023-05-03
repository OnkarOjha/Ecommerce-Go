package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/order"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func MakePaymentHandler(context *gin.Context) {
	utils.SetHeader(context)

	var paymentRequest request.OrderRequest

	utils.RequestDecoding(context, &paymentRequest)

	err := validation.CheckValidation(&paymentRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	order.MakePaymentService(context, paymentRequest)
}

func GetOrderDetails(context *gin.Context) {
	utils.SetHeader(context)

	order.GetOrderDetails(context)
}
