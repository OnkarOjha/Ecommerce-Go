package handler

import (
	"github.com/gin-gonic/gin"
	"main/server/context"
	"main/server/model"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"
	"main/server/validation"
)

func UserRegisterHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var registerRequest context.UserRequest

	utils.RequestDecoding(ctx, &registerRequest)

	err := validation.CheckValidation(&registerRequest)

	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.RegisterUserService(ctx, registerRequest)

}

func UserLoginHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var userLogin context.UserLogin

	err := utils.RequestDecoding(ctx, &userLogin)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&userLogin)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.UserLoginService(ctx, userLogin)
}

func UserVerifyOtpHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var verifyOtpRequest context.VerifyOtp

	err := utils.RequestDecoding(ctx, &verifyOtpRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&verifyOtpRequest)

	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.UserVerifyService(ctx, verifyOtpRequest)
}

func GetUserByIdHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	user.GetUserByIdService(ctx)

}

func EditUserProfile(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var editUserRequest context.EditUser

	err := utils.RequestDecoding(ctx, &editUserRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&editUserRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.EditUserService(ctx, editUserRequest)

}

func UserLogoutHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	user.LogoutUserService(ctx)

}

func UserAddressHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var userAddress model.UserAddresses

	err := utils.RequestDecoding(ctx, &userAddress)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = userAddress.Validate()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.UserAddressService(ctx, userAddress)

}

func UserAddressRetrieveHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	user.UserAddressRetrieveService(ctx)
}
