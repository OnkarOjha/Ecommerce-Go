package handler

import (
	"main/server/context"
	"main/server/model"
	"main/server/response"
	"main/server/services/user"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

// @Summary  		User Register Handler
// @Description  	Registering User with initial details in DB
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		user-register body string true "name and phone number of the user" SchemaExample({  "username" : "string", "usercontact" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/register [post]
func UserRegisterHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var registerRequest context.UserRequest

	err := utils.RequestDecoding(ctx, &registerRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to decode request")
		return
	}

	err = validation.CheckValidation(&registerRequest)

	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	user.RegisterUserService(ctx, registerRequest)

}

// @Summary  		User Login Handler
// @Description  	Login User with phone number
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		user-login body string true "phone number of the user" SchemaExample({ "usercontact" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/login [post]
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

// @Summary  		User Verify OTP Handler
// @Description  	Verify the OTP against the provided phone number
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		verify-otp body string true "phone number and otp of the user" SchemaExample({ "usercontact" : "string" , "otp" : "string" })
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/verify-otp [post]
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

// @Summary  		User Profile Details
// @Description  	This Handler provides all the user information with ID from token
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/get-user [get]
func GetUserByIdHandler(ctx *gin.Context) {

	utils.SetHeader(ctx)

	user.GetUserByIdService(ctx)

}

// @Summary  		User Edit Profile Details
// @Description  	This Handler enables user to edit his/her details
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		edit-profile body string true "name and gender of user" SchemaExample({ "usercontact" : "string" , "gender" : "string" })
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/edit-user [post]
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

// @Summary  		User Logout
// @Description  	This Handler will Log out the user
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/logout [delete]
func UserLogoutHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	user.LogoutUserService(ctx)

}

// @Summary  		User Address
// @Description  	This Handler will set user address
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		user-address-set body string true "user address details" SchemaExample({  "name" : "string","street" : "string","city" : "string","state" : "string","postalCode" : "string","country" : "string","phone" : "string","email" : "string","addressType" : "home/work/default" })
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/user/address [post]
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

// @Summary  		User Address Retrieve Handler
// @Description  	This Handler will get user addresses by passing "addresstype" param as "WORK/DEFAULT/HOME"
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param			addresstype query string true "address type details" SchemaExample({"addresstype" : "WORK/DEFAULT/HOME"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/user/address-get [get]
func UserAddressRetrieveHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	user.UserAddressRetrieveService(ctx)
}
