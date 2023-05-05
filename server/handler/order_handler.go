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

	err := utils.RequestDecoding(context, &paymentRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&paymentRequest)
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

func CancelOrderHandler(context *gin.Context) {
	utils.SetHeader(context)

	var cancelOrderRequest request.CancelOrderRequest
	err := utils.RequestDecoding(context, &cancelOrderRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&cancelOrderRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	order.CancelOrderService(context, cancelOrderRequest)
}

func MakeCartPaymentHandler(context *gin.Context) {
	utils.SetHeader(context)
	var paymentRequest request.CartOrderRequest
	err := utils.RequestDecoding(context, &paymentRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&paymentRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}
	order.MakeCartPaymentService(context, paymentRequest)
}
