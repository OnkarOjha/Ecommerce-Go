package handler

import (
	"main/server/context"
	"main/server/model"
	"main/server/response"
	"main/server/services/product"
	"main/server/services/vendors"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary  		Vendor Register Handler
// @Description  	Registering Vendor with initial details in DB
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		vendor-register body string true "vendor details" SchemaExample({  "gstNumber" : "29ABCDE1234F1Z5","companyName" : "Sports Tak","companyContact" : "9877370350","street" : "saytan gali kholi number 420","city" : "mohali","state" : "punjab","postalCode" : "152001","country" : "india"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-register [post]
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

	vendors.VendorRegisterService(ctx, vendorRegisterRequest)
}

// @Summary  		Vendor Login Handler
// @Description  	Login Vendor with initial details in DB
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		vendor-register body string true "vendor details" SchemaExample({  "gstNumber" : "29ABCDE1234F1Z5","companyContact" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-login [post]
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

	vendors.VendorLoginService(ctx, vendorLoginRequest)
}

// @Summary  		Vendor Verify OTP Handler
// @Description  	Verify the OTP against the provided phone number
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		verify-otp body string true "phone number and otp of the user" SchemaExample({ "contactNumber" : "string" , "otp" : "string" })
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-verify-otp [post]
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

	vendors.VerifyOtpService(ctx, verifyOtpRequest)
}

// @Summary  		Vendor Logout
// @Description  	This Handler will Log out the user
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-logout [delete]
func VendorLogoutHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	vendors.VendorLogoutService(ctx)

}

// @Summary  		Vendor Edit Profile Details
// @Description  	This Handler enables Vendor to edit his/her details
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		edit-profile body string true "company details" SchemaExample({  "companyName": "Ambani seth","description" : "Abra ka dabra jbewijbwr","logo" : "/home/chicmic/Downloads/test.jpg","bannerImage" : "/home/chicmic/Downloads/test.jpg"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-edit-details [post]
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

	vendors.VendorEditDetailsService(ctx, editDetailsRequest)

}

// @Summary  		Product Add from vendor side
// @Description  	This Handler adds multiple products from vendor side
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		product-add body string true "product description" SchemaExample({   "productName" : "Redmi Note 10 pro 4G","productDescription" : "4G smart phone","productPrice" : 9999.9,"productBrand" : "Redmi","productCategory" : "mobile","productImageUrl" : "http://dummyimage.com/169x100.png/cc0000/ffffff"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-product-add [post]
func InventoryProductAddHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventory model.Products

	err := utils.RequestDecoding(ctx, &productInventory)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = productInventory.ValidateProduct()
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductAddService(ctx, productInventory)

}

// @Summary  		Product Inventory Update
// @Description  	This Handler Update Product Inventory Details
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		product-inventory-update body string true "product id and product inventory" SchemaExample({  "productId" : "string","productInventory" : "float64"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-product-update [post]
func InventoryProductUpdateHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventoryEdit model.Products

	err := utils.RequestDecoding(ctx, &productInventoryEdit)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductUpdateService(ctx, productInventoryEdit)
}

// @Summary  		Product Delete from Inventory
// @Description  	This Handler will delete product from the inventory
// @Tags 			Vendor
// @Accept 			json
// @Procedure 		json
// @Param   		inventory-delete body string true "product id" SchemaExample({  "productId" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/vendor-product-delete [post]
func InventoryProductDeleteHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var productInventoryDelete context.ProductDeleteRequest

	err := utils.RequestDecoding(ctx, &productInventoryDelete)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	product.InventoryProductDeleteService(ctx, productInventoryDelete)

}

func VendorFileUpload(ctx *gin.Context) {
	utils.SetHeader(ctx)

	vendors.FileUpload(ctx)
}

func VendorFileGet(ctx *gin.Context) {
	utils.SetHeader(ctx)

	vendors.FileGet(ctx)
}
