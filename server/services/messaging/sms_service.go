package messaging

import (
	"errors"
	"fmt"
	"main/server/services/twilio"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/crypto/bcrypt"
)

//Generate random OTP
func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(9999-1111) + 1000
	return otp
}

//Send SMS to given number
//Sets two cookies OTP and NUMBER which contains values in hashed form
func SendSmsService(ctx *gin.Context, number string) (bool, *string) {
	otp := GenerateOTP()
	params := &api.CreateMessageParams{}
	params.SetBody(fmt.Sprint(otp))
	params.SetFrom("+15076046305")
	params.SetTo("+91" + number)
	resp, err := twilio.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		return false, nil
	} else {
		str := strconv.Itoa(otp)
		otpHash, _ := bcrypt.GenerateFromPassword([]byte(str), 8)
		otpHashed := string(otpHash)
		otpCookie := &http.Cookie{
			Name:   "otp",
			Value:  otpHashed,
			MaxAge: 60,
		}
		http.SetCookie(ctx.Writer, otpCookie)

		numberHash, _ := bcrypt.GenerateFromPassword([]byte(number), 8)
		numberHashed := string(numberHash)
		numberCookie := &http.Cookie{
			Name:   "number",
			Value:  numberHashed,
			MaxAge: 60,
		}
		http.SetCookie(ctx.Writer, numberCookie)
		return true, resp.Sid
	}
}

// Check OTP and number from cookies
func CheckOtpService(ctx *gin.Context, originalNumber string, numberProvided string, originalOtp string, otpProvided string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(originalNumber), []byte(numberProvided))
	if err != nil {
		fmt.Println("dkngskd", err.Error())
		return false, errors.New("Error comparing number with hashed number")
	}
	err = bcrypt.CompareHashAndPassword([]byte(originalOtp), []byte(otpProvided))
	if err != nil {

		return false, errors.New("Error comparing otp with hashed otp")
	}

	return true, nil
}
