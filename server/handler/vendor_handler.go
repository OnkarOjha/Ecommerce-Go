package handler

import (
	"main/server/context"
	"main/server/response"
	"main/server/services/vendor"
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

	vendor.VendorRegisterService(ctx, vendorRegisterRequest)
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

	vendor.VendorLoginService(ctx, vendorLoginRequest)
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

	vendor.VerifyOtpService(ctx, verifyOtpRequest)
}

func VendorLogoutHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	vendor.VendorLogoutService(ctx)

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

	vendor.VendorEditDetailsService(ctx, editDetailsRequest)

}
