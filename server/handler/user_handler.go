package handler

import (
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(context *gin.Context) {
	utils.SetHeader(context)

	var registerRequest request.UserRequest

	utils.RequestDecoding(context, &registerRequest)

	err := validation.CheckValidation(&registerRequest)

	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.RegisterUserService(context, registerRequest)

}

func UserLoginHandler(context *gin.Context) {
	utils.SetHeader(context)

	var userLogin request.UserLogin

	err := utils.RequestDecoding(context, &userLogin)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&userLogin)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.UserLoginService(context, userLogin)
}

func UserVerifyOtpHandler(context *gin.Context) {
	utils.SetHeader(context)

	var verifyOtpRequest request.VerifyOtp

	err := utils.RequestDecoding(context, &verifyOtpRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&verifyOtpRequest)

	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.VerifyOtpService(context, verifyOtpRequest)
}

func GetUserByIdHandler(context *gin.Context) {

	utils.SetHeader(context)

	var getUser request.GetUser

	err := utils.RequestDecoding(context, &getUser)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&getUser)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.GetUserByIdService(context, getUser.UserId)

}

func EditUserProfile(context *gin.Context) {
	utils.SetHeader(context)

	var editUserRequest request.EditUser

	err := utils.RequestDecoding(context, &editUserRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&editUserRequest)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.EditUserService(context, editUserRequest)

}

func UserLogoutHandler(context *gin.Context) {
	utils.SetHeader(context)

	var logoutUser request.LogoutUser

	err := utils.RequestDecoding(context, &logoutUser)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = validation.CheckValidation(&logoutUser)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.LogoutUserService(context, logoutUser)

}

func UserAddressHandler(context *gin.Context) {
	utils.SetHeader(context)

	var userAddress model.UserAddresses

	err := utils.RequestDecoding(context, &userAddress)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	err = userAddress.Validate()
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.UserAddressService(context, userAddress)

}

func UserAddressRetrieveHandler(context *gin.Context) {
	utils.SetHeader(context)

	user.UserAddressRetrieveService(context)
}
