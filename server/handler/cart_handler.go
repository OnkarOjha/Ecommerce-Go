package handler

import (
	"main/server/context"
	"main/server/response"
	"main/server/services/cart"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func AddToCartHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)
	var addToCartRequest context.AddToCartRequest

	err := utils.RequestDecoding(ctx, &addToCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&addToCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	cart.AddToCartService(ctx, addToCartRequest)

}

func AddProductHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var addToCartRequest context.AddToCartRequest

	err := utils.RequestDecoding(ctx, &addToCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&addToCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	cart.AddProductService(ctx, addToCartRequest)

}

func RemoveFromCartHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var removeFromCartRequest context.RemoveFromCart

	err := utils.RequestDecoding(ctx, &removeFromCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&removeFromCartRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	cart.RemoveFromCartService(ctx, removeFromCartRequest)
}

func RemoveProductHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var removeProductFromCart context.RemoveProduct

	err := utils.RequestDecoding(ctx, &removeProductFromCart)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&removeProductFromCart)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	cart.RemoveProductService(ctx, removeProductFromCart)

}

func GetCartDetailsHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)
	cart.GetCartDetailsService(ctx)
}
