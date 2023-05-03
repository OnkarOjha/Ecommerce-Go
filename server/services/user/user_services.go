package user

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/provider"
	"main/server/request"
	"main/server/response"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var twilioClient *twilio.RestClient

func TwilioInit(password string) {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: password,
	})
}

func RegisterUserService(context *gin.Context, registerRequest request.UserRequest) {
	// check if user already registered
	if db.RecordExist("users", "contact", registerRequest.UserContact) {
		response.ErrorResponse(context, 403, "User already registered , please proceed to login")
		return
	}

	var user model.User
	user.Contact = registerRequest.UserContact
	user.UserName = registerRequest.UserName

	err := db.CreateRecord(&user)
	if err != nil {
		response.ErrorResponse(context, 500, err.Error())
		return
	}
	response.ShowResponse(
		"Success",
		200,
		"User Registered successfully , please proceed to login",
		user,
		context,
	)
}

func UserLoginService(context *gin.Context, userLogin request.UserLogin) {
	// check that number is registered or not
	if !db.RecordExist("users", "contact", userLogin.UserContact) {
		response.ErrorResponse(context, 400, "Number not registered , please register")
		return
	} else {
		if db.RecordExist("users", "is_active", "false") {
			fmt.Println("signing in  old user")
			// it's an old user who is trying to login
			fmt.Println(twilioClient.Client.AccountSid())
			ok, sid := SendOtpService(context, "+91"+userLogin.UserContact)
			if ok {
				response.ShowResponse("Success", 200, "Welcome back! OTP sent successfully", sid, context)
				return
			} else {
				response.ErrorResponse(context, 400, "Error sending OTP")
				return
			}

		}
	}
	// new user logging in
	ok, sid := SendOtpService(context, "+91"+userLogin.UserContact)
	if ok {
		response.ShowResponse("Success", 200, "OTP sent successfully", sid, context)
		return
	} else {
		response.ErrorResponse(context, 400, "Error sending OTP")
		return
	}

}

func SendOtpService(context *gin.Context, contact string) (bool, *string) {
	params := &openapi.CreateVerificationParams{}
	fmt.Println("contact", contact)
	params.SetTo(contact)

	params.SetChannel("sms")
	fmt.Println("veify service sid: ", os.Getenv("VERIFY_SERVICE_SID"))
	resp, err := twilioClient.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		return false, nil
	} else {
		return true, resp.Sid
	}
}

func VerifyOtpService(context *gin.Context, verifyOtpRequest request.VerifyOtp) {
	params := &openapi.CreateVerificationCheckParams{}

	params.SetTo("+91" + verifyOtpRequest.UserContact)

	params.SetCode(verifyOtpRequest.Otp)
	fmt.Println("params:", verifyOtpRequest)
	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		response.ErrorResponse(context, 401, "Verification Failed")
		return
	} else if *resp.Status == "approved" {
		// user creation
		var tokenClaims model.Claims
		var user model.User
		var userSession model.Session
		err := db.FindById(&user, verifyOtpRequest.UserContact, "contact")
		fmt.Println("user: ", user)
		if err != nil {
			response.ErrorResponse(context, 500, "Error finding user in DB")
			return
		}
		user.Is_Active = true
		tokenClaims.UserId = user.UserId
		tokenClaims.Phone = user.Contact
		tokenClaims.Role = "customer"
		tokenClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 2))

		db.UpdateRecord(&user, user.UserId, "user_id")

		tokenString := provider.GenerateToken(tokenClaims, context)

		userSession.Token = tokenString
		userSession.UserId = user.UserId

		err = db.CreateRecord(&userSession)
		if err != nil {
			response.ErrorResponse(context, 500, "Error creating record: "+err.Error())
			return
		}
		response.ShowResponse("Success", 200, "User verified successfully", user, context)
		response.ShowResponse("Success", 200, "Session created successfully", userSession, context)

	} else {
		response.ErrorResponse(context, 401, "Verification Failed")
		return
	}
}

func GetUserByIdService(context *gin.Context, userId string) {
	if !db.RecordExist("users", "user_id", userId) {
		response.ErrorResponse(context, 400, "User not found")
		return
	}

	var user model.User
	err := db.FindById(&user, userId, "user_id")

	if err != nil {
		response.ErrorResponse(context, 400, "User not found")
	}

	response.ShowResponse("Success", 200, "User Fetched successfully", user, context)
}

func EditUserService(context *gin.Context, editUserRequest request.EditUser) {
	if !db.RecordExist("users", "user_id", editUserRequest.UserId) {
		response.ErrorResponse(context, 400, "User not found")
		return
	}
	var user model.User
	err := db.FindById(&user, editUserRequest.UserId,
		"user_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error finding user")
	}

	user.UserName = editUserRequest.UserName
	user.Gender = editUserRequest.Gender
	db.UpdateRecord(&user, editUserRequest.UserId, "user_id")

	response.ShowResponse("Success", 200, "User Profile updated successfully", user, context)
}

func LogoutUserService(context *gin.Context, logoutUser request.LogoutUser) {
	var user model.User
	err := db.FindById(&user, logoutUser.UserId,
		"user_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error finding user")
		return
	}

	user.Is_Active = false
	query := "UPDATE users set is_active=false where user_id=?"
	db.QueryExecutor(query , user , user.UserId)

	var userSession  model.Session
	err = db.FindById(&userSession , logoutUser.UserId ,"user_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error finding user session")
		return
	}
	err = db.DeleteRecord(&userSession , userSession.
	UserId , "user_id")
	if err!=nil{
		response.ErrorResponse(context , 400, "Error deleting user session")
		return
	}

	response.ShowResponse("Success" , 200 , "Logout Successfull" , user , context)
}
