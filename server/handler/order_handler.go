package handler

import (
	"main/server/context"
	"main/server/response"
	"main/server/services/order"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func MakePaymentHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var paymentRequest context.OrderRequest

	err := utils.RequestDecoding(ctx, &paymentRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&paymentRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	order.MakePaymentService(ctx, paymentRequest)
}

func GetOrderDetails(ctx *gin.Context) {
	utils.SetHeader(ctx)

	order.GetOrderDetails(ctx)
}

func CancelOrderHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var cancelOrderRequest context.CancelOrderRequest
	err := utils.RequestDecoding(ctx, &cancelOrderRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&cancelOrderRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	order.CancelOrderService(ctx, cancelOrderRequest)
}

func MakeCartPaymentHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)
	var paymentRequest context.CartOrderRequest
	err := utils.RequestDecoding(ctx, &paymentRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&paymentRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}
	order.MakeCartPaymentService(ctx, paymentRequest)
}
