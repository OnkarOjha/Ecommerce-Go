package vendors

import (
	"io/ioutil"
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/messaging"
	"main/server/services/order"
	"main/server/services/token"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Vendor Register service
func VendorRegisterService(ctx *gin.Context, vendorRegisterRequest context.VendorRegisterRequest) {
	// check if user already registered
	if db.RecordExist("vendors", "vendor_contact", vendorRegisterRequest.CompanyContact) {
		response.ErrorResponse(ctx, 403, "Vendor already registered , please proceed to login")
		return
	}
	if db.RecordExist("users", "contact", vendorRegisterRequest.CompanyContact) {
		response.ErrorResponse(ctx, 403, "Vendor already registered as a user , please provide a different number")
		return
	}

	var vendor model.Vendor

	vendor.VendorName = vendorRegisterRequest.CompanyName
	vendor.VendorContact = vendorRegisterRequest.CompanyContact
	vendor.Street = vendorRegisterRequest.Street
	vendor.City = vendorRegisterRequest.City
	vendor.State = vendorRegisterRequest.State
	vendor.PostalCode = vendorRegisterRequest.PostalCode
	vendor.Country = vendorRegisterRequest.Country
	vendor.GstIn = vendorRegisterRequest.GstNumber

	err := db.CreateRecord(&vendor)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
		return
	}
	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Vendor Registered successfully , please proceed to login",
		vendor,
		ctx,
	)
}

// Venor Login service
func VendorLoginService(ctx *gin.Context, vendorLoginRequest context.VendorLoginRequest) {
	//check if the gst number is registered with the mobile number
	if db.BothExists("vendors", "gst_in", vendorLoginRequest.GstNumber, "vendor_contact", vendorLoginRequest.CompanyContact) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No vendor Registered with this GST Number and contact number")
		return
	}

	if db.RecordExist("vendors", "is_active", "false") {
		// old user logging in
		ok, sid := messaging.SendSmsService(ctx, vendorLoginRequest.CompanyContact)
		if ok {
			response.ShowResponse(
				"Success",
				utils.HTTP_OK,
				"Welcome back! OTP sent successfully",
				sid,
				ctx,
			)
			return
		} else {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error sending OTP")
			return
		}
	}
	// new user logging in
	ok, sid := messaging.SendSmsService(ctx, vendorLoginRequest.CompanyContact)
	if ok {
		response.ShowResponse("Success", utils.HTTP_OK, "OTP sent successfully", sid, ctx)
		return
	} else {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error sending OTP")
		return
	}
}

// Vendor OTP Sending service
func VerifyOtpService(ctx *gin.Context, verifyOtpRequest context.VendorVerifyOtpRequest) {
	otpCookie, err := ctx.Request.Cookie("otp")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}
	numberCookie, err := ctx.Request.Cookie("number")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	otpVerified, err := messaging.CheckOtpService(ctx, numberCookie.Value, verifyOtpRequest.ContactNumber, otpCookie.Value, verifyOtpRequest.Otp)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, err.Error())
		return
	}

	if otpVerified {
		var tokenClaims token.Claims
		var vendor model.Vendor
		var userSession model.Session
		err := db.FindById(&vendor, verifyOtpRequest.ContactNumber, "vendor_contact")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user in DB")
			return
		}
		vendor.Is_Active = true
		tokenClaims.UserId = vendor.VendorId
		tokenClaims.Phone = vendor.VendorContact
		tokenClaims.Role = "vendor"
		tokenClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * time.Hour * 24))

		db.UpdateRecord(&vendor, vendor.VendorId, "vendor_id")

		tokenString := token.GenerateToken(tokenClaims, ctx)

		userSession.Token = tokenString
		userSession.UserId = vendor.VendorId

		err = db.CreateRecord(&userSession)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
			return
		}
		response.ShowResponse("Success", utils.HTTP_OK, "Vendor verified successfully", vendor, ctx)
		response.ShowResponse("Success", utils.HTTP_OK, "Session created successfully", userSession, ctx)

	} else {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Verification Failed")
		return
	}

}

// Vendor logout and session delete
func VendorLogoutService(ctx *gin.Context) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}

	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}
	var vendor model.Vendor
	err = db.FindById(&vendor, vendorId,
		"vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor")
		return
	}

	vendor.Is_Active = false
	query := "UPDATE users set is_active=false where vendor_id=?"
	err = db.QueryExecutor(query, vendor, vendor.VendorId)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "unable to execute query")
		return
	}

	var userSession model.Session
	err = db.FindById(&userSession, vendorId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor session")
		return
	}
	err = db.DeleteRecord(&userSession, userSession.
		UserId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error deleting vendor session")
		return
	}

	response.ShowResponse("Success", utils.HTTP_OK, "Logout Successfull", vendor, ctx)
}

// Edit vendor details
func VendorEditDetailsService(ctx *gin.Context, editDetailsRequest context.VendorEditDetailsRequest) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}

	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}

	var vendor model.Vendor
	err = db.FindById(&vendor, vendorId,
		"vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor")
		return
	}

	vendor.GstIn = editDetailsRequest.GstNumber
	vendor.Street = editDetailsRequest.Street
	vendor.City = editDetailsRequest.City
	vendor.State = editDetailsRequest.State
	vendor.PostalCode = editDetailsRequest.PostalCode
	vendor.Country = editDetailsRequest.Country
	vendor.Description = editDetailsRequest.Description
	// vendor.BannerImage = bannerContent
	// vendor.Logo = imageContent

	db.UpdateRecord(&vendor, vendorId, "vendor_id")

	err = db.FindById(&vendor, vendorId,
		"vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor")
		return
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Edit Details Successful",
		vendor,
		ctx,
	)

}

func FileUpload(ctx *gin.Context) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}

	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}

	// Get the uploaded file
	bannerImage, _, err := ctx.Request.FormFile("bannerimage")
	if err != nil {

		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding file")
		return
	}
	defer bannerImage.Close()

	logoImage, _, err := ctx.Request.FormFile("logoimage")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding file")
		return
	}
	defer logoImage.Close()

	// Read the file content
	bannerContent, err := ioutil.ReadAll(bannerImage)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error reading file content")
		return
	}

	imageContent, err := ioutil.ReadAll(logoImage)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error reading file content")
		return
	}

	var vendor model.Vendor

	// fmt.Println("vendfor", vendor)
	err = db.FindById(&vendor, vendorId,
		"vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor")
		return
	}
	vendor.BannerImage = bannerContent
	vendor.Logo = imageContent
	db.UpdateRecord(&vendor, vendorId, "vendor_id")

	err = db.FindById(&vendor, vendorId,
		"vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding vendor")
		return
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Upload file successfull",
		vendor,
		ctx,
	)

}

func FileGet(ctx *gin.Context) {
	// vendorId, err := order.IdFromToken(ctx)
	// if err != nil {
	// 	response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
	// 	return
	// }

	// if vendorId == "" {
	// 	response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
	// 	return
	// }
	// err = os.MkdirAll("./File/uploads", os.ModePerm)
	// if err != nil {
	// 	response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Failed initialising directory upload")
	// 	return

	// }

	// var filecontent []byte
	// query := "select logo from vendors where vendor_id='" + vendorId + "'"
	// db.QueryExecutor(query, &filecontent)

	// dst, err := os.Create(fmt.Sprintf("./File/uploads/%d", time.Now().UnixNano()))

	// if err != nil {
	// 	response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Failed assigning file name")
	// 	return
	// }

	// defer dst.Close()
	// // Write the file to disk
	// // err = ioutil.WriteFile(dst.Name(), filecontent, 0644)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// n, err := dst.Write(filecontent)
	// if err != nil {
	// 	response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Failed writing to file")
	// 	return
	// }
	// response.ShowResponse(
	// 	"Success",
	// 	200,
	// 	"File length displayed",
	// 	n,
	// 	ctx,
	// )

	// // Set the Content-Disposition header to indicate the filename
	// ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))

	// // Set the Content-Type header based on the file extension
	// switch filepath.Ext(file.Filename) {
	// case ".pdf":
	// 	ctx.Header("Content-Type", "application/pdf")
	// case ".jpg", ".jpeg":
	// 	ctx.Header("Content-Type", "image/jpeg")
	// case ".png":
	// 	ctx.Header("Content-Type", "image/png")
	// default:
	// 	ctx.Header("Content-Type", "application/octet-stream")
	// }

	// // Stream the file content to the HTTP response
	// ctx.Status(http.StatusOK)
	// io.Copy(ctx.Writer, bytes.NewReader(file.Content))

}
