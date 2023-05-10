package handler

import (
	"main/server/context"
	"main/server/response"
	"main/server/services/cart"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

// @Summary  		Add to cart handler
// @Description  	Add products to cart
// @Tags 			Cart
// @Accept 			json
// @Procedure 		json
// @Param   		add-to-cart body string true "product id and product count" SchemaExample({  "productId" : "string", "productCount" : "float64"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/add-to-cart [post]
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

// @Summary  		Add Products
// @Description  	Add more products of same type in cart
// @Tags 			Cart
// @Accept 			json
// @Procedure 		json
// @Param   		add-product body string true "product id and product count" SchemaExample({  "productId" : "string", "productCount" : "float64"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/add-product [put]
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

// @Summary  		Remove product from cart
// @Description  	Remove the product totally from cart
// @Tags 			Cart
// @Accept 			json
// @Procedure 		json
// @Param   		remove-product body string true "CartId and ProductId" SchemaExample({  "cartId" : "string", "productId" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/remove-from-cart [delete]
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

// @Summary  		Reduce product Count
// @Description  	Reduce the product count from product
// @Tags 			Cart
// @Accept 			json
// @Procedure 		json
// @Param   		remove-product body string true "CartId , ProductId and product count" SchemaExample({  "cartId" : "string", "productId" : "string", "productCount" : "float64"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/remove-product [delete]
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

// @Summary  		Get Cart Details
// @Description  	Get the Cart Summary
// @Tags 			Cart
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/get-cart-details [get]
func GetCartDetailsHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)
	cart.GetCartDetailsService(ctx)
}
