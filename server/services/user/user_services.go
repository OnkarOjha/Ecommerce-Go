package user

import (
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/order"
	"main/server/services/token"
	"main/server/services/twilio"
	"main/server/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Register User Service
func RegisterUserService(ctx *gin.Context, registerRequest context.UserRequest) {
	// check if user already registered
	if db.RecordExist("users", "contact", registerRequest.UserContact) {
		response.ErrorResponse(ctx, 403, "User already registered , please proceed to login")
		return
	}

	var user model.User

	user.Contact = registerRequest.UserContact
	user.UserName = registerRequest.UserName

	err := db.CreateRecord(&user)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
		return
	}
	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"User Registered successfully , please proceed to login",
		user,
		ctx,
	)
}

// User Login Service
func UserLoginService(ctx *gin.Context, userLogin context.UserLogin) {
	// check that number is registered or not
	if !db.RecordExist("users", "contact", userLogin.UserContact) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Number not registered , please register")
		return
	} else {
		if db.RecordExist("users", "is_active", "false") {
			// it's an old user who is trying to login
			ok, sid := twilio.SendOtpService(ctx, "+91"+userLogin.UserContact)
			if ok {
				response.ShowResponse("Success", utils.HTTP_OK, "Welcome back! OTP sent successfully", sid, ctx)
				return
			} else {
				response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error sending OTP")
				return
			}

		}
	}

	// new user logging in
	ok, sid := twilio.SendOtpService(ctx, "+91"+userLogin.UserContact)
	if ok {
		response.ShowResponse("Success", utils.HTTP_OK, "OTP sent successfully", sid, ctx)
		return
	} else {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error sending OTP")
		return
	}
}

// User Verify Service with OTP and user signin
func UserVerifyService(ctx *gin.Context, verifyOtpRequest context.VerifyOtp) {
	veriyStatus, err := twilio.VerifyOtpService(ctx, verifyOtpRequest.UserContact, verifyOtpRequest.Otp)
	if veriyStatus == "approved" {
		// user creation
		var tokenClaims token.Claims
		var user model.User
		var userSession model.Session
		err := db.FindById(&user, verifyOtpRequest.UserContact, "contact")

		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user in DB")
			return
		}
		user.Is_Active = true
		tokenClaims.UserId = user.UserId
		tokenClaims.Phone = user.Contact
		tokenClaims.Role = "customer"
		tokenClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))

		db.UpdateRecord(&user, user.UserId, "user_id")

		tokenString := token.GenerateToken(tokenClaims, ctx)

		userSession.Token = tokenString
		userSession.UserId = user.UserId

		err = db.CreateRecord(&userSession)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
			return
		}
		response.ShowResponse("Success", utils.HTTP_OK, "User verified successfully", user, ctx)
		response.ShowResponse("Success", utils.HTTP_OK, "Session created successfully", userSession, ctx)

	}

	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Verification Failed")
		return
	}
}

// User Get service
func GetUserByIdService(context *gin.Context) {

	userId, err := order.IdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, utils.HTTP_BAD_REQUEST, "User Not Found")
		return
	}

	if !db.RecordExist("users", "user_id", userId) {
		response.ErrorResponse(context, utils.HTTP_BAD_REQUEST, "User not found")
		return
	}

	var user model.User
	err = db.FindById(&user, userId, "user_id")

	if err != nil {
		response.ErrorResponse(context, utils.HTTP_BAD_REQUEST, "User not found")
	}

	response.ShowResponse("Success", utils.HTTP_OK, "User Fetched successfully", user, context)
}

// Edit User Details
func EditUserService(ctx *gin.Context, editUserRequest context.EditUser) {
	userId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "User Not Found")
		return
	}

	if !db.RecordExist("users", "user_id", userId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "User not found")
		return
	}
	var user model.User
	err = db.FindById(&user, userId,
		"user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding user")
	}

	user.UserName = editUserRequest.UserName
	user.Gender = editUserRequest.Gender
	db.UpdateRecord(&user, userId, "user_id")

	response.ShowResponse("Success", utils.HTTP_OK, "User Profile updated successfully", user, ctx)
}

// User logout Service
func LogoutUserService(ctx *gin.Context) {
	userId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}

	if userId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No User ID provided")
		return
	}
	var user model.User
	err = db.FindById(&user, userId,
		"user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding user")
		return
	}

	user.Is_Active = false
	query := "UPDATE users set is_active=false where user_id=?"
	err = db.QueryExecutor(query, user, user.UserId)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to execute query")
		return
	}

	var userSession model.Session
	err = db.FindById(&userSession, userId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding user session")
		return
	}
	err = db.DeleteRecord(&userSession, userSession.
		UserId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error deleting user session")
		return
	}

	response.ShowResponse("Success", utils.HTTP_OK, "Logout Successfull", user, ctx)
}

// User Address Set DEFAULT , HOME , WORK
func UserAddressService(ctx *gin.Context, userAddressRequest model.UserAddresses) {

	userId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}

	if userId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No User ID provided")
		return
	}
	var dbConstants model.DbConstant
	dbConstants.ConstantName = strings.ToLower(userAddressRequest.AddressType)
	if strings.ToLower(userAddressRequest.AddressType) == "default" {
		dbConstants.ConstantShortHand = "DEFAULT"
	}
	if strings.ToLower(userAddressRequest.AddressType) == "home" {
		dbConstants.ConstantShortHand = "HOME"
	}
	if strings.ToLower(userAddressRequest.AddressType) == "work" {
		dbConstants.ConstantShortHand = "WORK"
	}
	userAddressRequest.AddressType = dbConstants.ConstantShortHand
	userAddressRequest.UserId = userId
	err = db.CreateRecord(&userAddressRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating DB record")
		return
	}

	DBConstantService(ctx, dbConstants)

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"User Address logged Successfully",
		userAddressRequest,
		ctx,
	)
}

// db constant set service
func DBConstantService(ctx *gin.Context, dbConstants model.DbConstant) {
	if dbConstants.ConstantShortHand == "DEFAULT" || dbConstants.ConstantShortHand == "HOME" || dbConstants.ConstantShortHand != "WORK" {
		exists1 := db.RecordExist("db_constants", "constant_short_hand", "DEFAULT")
		exists2 := db.RecordExist("db_constants", "constant_short_hand", "HOME")
		exists3 := db.RecordExist("db_constants", "constant_short_hand", "WORK")
		if !exists1 || !exists2 || !exists3 {
			err := db.CreateRecord(&dbConstants)
			if err != nil {
				response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating DB record")
				return
			}
		}
	}
}

// user address retrieve
func UserAddressRetrieveService(ctx *gin.Context) {
	addressType := ctx.Query("addresstype")

	if addressType == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Address type not specified")
		return
	}
	var userAddress []model.UserAddresses
	if addressType == "DEFAULT" {
		query := "SELECT * FROM user_addresses WHERE address_type='" + addressType + "'"
		err := db.QueryExecutor(query, &userAddress)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user address")
			return
		}
	}
	if addressType == "HOME" {
		query := "SELECT * FROM user_addresses WHERE address_type='" + addressType + "'"
		err := db.QueryExecutor(query, &userAddress)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user address")
			return
		}
	}
	if addressType == "WORK" {
		query := "SELECT * FROM user_addresses WHERE address_type='" + addressType + "'"
		err := db.QueryExecutor(query, &userAddress)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user address")
			return
		}
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"User Address retrieved successfully",
		userAddress,
		ctx,
	)
}
