package handler

import (
	"main/server/request"
	"main/server/response"
	"main/server/services/cart"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func AddToCartHandler(context *gin.Context) {
	utils.SetHeader(context)

	var addToCartRequest request.AddToCartRequest

	err := utils.RequestDecoding(context, &addToCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&addToCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	cart.AddToCartService(context, addToCartRequest)

}

func AddProductHandler(context *gin.Context) {
	utils.SetHeader(context)

	var addToCartRequest request.AddToCartRequest

	err := utils.RequestDecoding(context, &addToCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&addToCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	cart.AddProductService(context, addToCartRequest)

}

func RemoveFromCartHandler(context *gin.Context) {
	utils.SetHeader(context)

	var removeFromCartRequest request.RemoveFromCart

	err := utils.RequestDecoding(context, &removeFromCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&removeFromCartRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	cart.RemoveFromCartService(context, removeFromCartRequest)
}

func RemoveProductHandler(context *gin.Context) {
	utils.SetHeader(context)

	var removeProductFromCart request.RemoveProduct

	err := utils.RequestDecoding(context, &removeProductFromCart)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&removeProductFromCart)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	cart.RemoveProductService(context, removeProductFromCart)

}

func GetCartDetailsHandler(context *gin.Context) {
	utils.SetHeader(context)
	cart.GetCartDetailsService(context)
}
