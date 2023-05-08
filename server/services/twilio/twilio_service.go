package twilio

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var TwilioClient *twilio.RestClient

func TwilioInit(password string) {
	TwilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: password,
	})
}

func SendOtpService(ctx *gin.Context, contact string) (bool, *string) {
	params := &openapi.CreateVerificationParams{}

	params.SetTo(contact)

	params.SetChannel("sms")

	resp, err := TwilioClient.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		return false, nil
	} else {
		return true, resp.Sid
	}
}

func VerifyOtpService(ctx *gin.Context, number string, otp string) (string, error) {
	params := &openapi.CreateVerificationCheckParams{}

	params.SetTo("+91" + number)

	params.SetCode(otp)

	resp, err := TwilioClient.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		return "", err
	} else if *resp.Status == "approved" {
		return *resp.Status, err
	}

	return "", nil

}
