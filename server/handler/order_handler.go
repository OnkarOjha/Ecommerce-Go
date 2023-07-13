package handler

import (
	"main/server/context"
	"main/server/response"
	"main/server/services/order"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

// @Summary  		Make Payment Handler
// @Description  	Make payment for products from cart
// @Tags 			Payment
// @Accept 			json
// @Procedure 		json
// @Param   		payment body string true "Order Payment" SchemaExample({   "productId" : "string","cartId" : "string","cardNumber": "string","expMonth": "string","expYear": "string","cvc": "string","couponName" : "string")
// @Param			addresstype query string true "address type" SchemaExample({"addresstype" : "DEFAULT/WORK/HOME"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/payment [post]
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

// @Summary  		Get Order Details
// @Description  	Get the full order details made till now
// @Tags 			Payment
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/order-details [get]
func GetOrderDetails(ctx *gin.Context) {
	utils.SetHeader(ctx)

	order.GetOrderDetails(ctx)
}

// @Summary  		Cancel Order
// @Description  	Cancel Order
// @Tags 			Payment
// @Accept 			json
// @Procedure 		json
// @Param   		orderId body string true "Provide Order Id to cancel order" SchemaExample({  "orderId" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/cancel-order [put]
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

// @Summary  		Make Cart Payment Handler
// @Description  	Make payment for the whole cart
// @Tags 			Payment
// @Accept 			json
// @Procedure 		json
// @Param   		payment body string true "Cart Order Payment" SchemaExample({   "productId" : "string","cartId" : "string","cardNumber": "string","expMonth": "string","expYear": "string","cvc": "string","couponName" : "string")
// @Param			addresstype query string true "address type" SchemaExample({"addresstype" : "DEFAULT/WORK/HOME"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/cart-payment [post]
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

// @Summary  		Vendor Order Status
// @Description  	Vendor Can update order status according to dispatched/confirmed
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		orderId body string true "Provide Order Id to cancel order" SchemaExample({  "orderId" : "string" ,"orderStatus" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-product-status [post]
func VendorOrderStatusUpdateHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var orderUpdateRequest context.VendorOrderStatusUpdate

	err := utils.RequestDecoding(ctx, &orderUpdateRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&orderUpdateRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	order.VendorOrderStatusUpdateService(ctx, orderUpdateRequest)
}
