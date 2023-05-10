package handler

import (
	"main/server/context"
	"main/server/model"
	"main/server/response"
	"main/server/services/product"
	"main/server/services/vendors"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func VendorRegisterHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var vendorRegisterRequest context.VendorRegisterRequest

	err := utils.RequestDecoding(ctx, &vendorRegisterRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = vendorRegisterRequest.ValidateRegister()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	vendors.VendorRegisterService(ctx, vendorRegisterRequest)
}

func VendorLoginHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var vendorLoginRequest context.VendorLoginRequest

	err := utils.RequestDecoding(ctx, &vendorLoginRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}
	err = vendorLoginRequest.ValidateLogin()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	vendors.VendorLoginService(ctx, vendorLoginRequest)
}

func VenderVerifyOtpHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var verifyOtpRequest context.VendorVerifyOtpRequest

	err := utils.RequestDecoding(ctx, &verifyOtpRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = verifyOtpRequest.ValidateOtp()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	vendors.VerifyOtpService(ctx, verifyOtpRequest)
}

func VendorLogoutHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	vendors.VendorLogoutService(ctx)

}

func VendorEditDetailsHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var editDetailsRequest context.VendorEditDetailsRequest

	err := utils.RequestDecoding(ctx, &editDetailsRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = editDetailsRequest.ValidateEditDetails()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	vendors.VendorEditDetailsService(ctx, editDetailsRequest)

}

func InventoryProductAddHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventory model.Products

	err := utils.RequestDecoding(ctx, &productInventory)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = productInventory.ValidateProduct()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductAddService(ctx, productInventory)

}

func InventoryProductUpdateHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventoryEdit model.Products

	err := utils.RequestDecoding(ctx, &productInventoryEdit)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductUpdateService(ctx, productInventoryEdit)
}

func InventoryProductDeleteHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventoryDelete context.ProductDeleteRequest

	err := utils.RequestDecoding(ctx, &productInventoryDelete)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductDeleteService(ctx, productInventoryDelete)

}
