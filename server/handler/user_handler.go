package handler

import (
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

	user.RegisterUserService(context ,registerRequest)

}

func UserLoginHandler(context *gin.Context){
	utils.SetHeader(context)

	var userLogin request.UserLogin

	utils.RequestDecoding(context, &userLogin)

	err := validation.CheckValidation(&userLogin)
	if err!=nil{
		response.ErrorResponse(context , 400 , err.Error())
		return
	}
	
	user.UserLoginService(context , userLogin)
}


func UserVerifyOtpHandler(context *gin.Context){
	utils.SetHeader(context)

	var verifyOtpRequest request.VerifyOtp

	utils.RequestDecoding(context , &verifyOtpRequest)

	err := validation.CheckValidation(&verifyOtpRequest)

	if err!=nil{
		response.ErrorResponse(context , 400 , err.Error())
		return
	}

	user.VerifyOtpService(context , verifyOtpRequest)
}

func GetUserByIdHandler(context *gin.Context) {

	utils.SetHeader(context)

	var getUser request.GetUser

	utils.RequestDecoding(context, &getUser)

	err := validation.CheckValidation(&getUser)
	if err != nil {
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.GetUserByIdService(context, getUser.UserId)

}

func EditUserProfile(context *gin.Context){
	utils.SetHeader(context)

	var editUserRequest request.EditUser

	utils.RequestDecoding(context , &editUserRequest)

	err := validation.CheckValidation(&editUserRequest)
	if err!=nil{
		response.ErrorResponse(context, 400, err.Error())
		return
	}

	user.EditUserService(context , editUserRequest)

}