package handler

import (
	"main/server/model"
	"main/server/response"
	"main/server/services/coupons"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

// @Summary  		Add Coupon Handler
// @Description  	Add Coupon Handler
// @Tags 			Coupon
// @Accept 			json
// @Procedure 		json
// @Param   		add-coupon body string true "Coupon Name and Coupon price" SchemaExample({  "couponName" : "string", "couponPrice" : "float64"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/coupon-add [post]
func AddCouponHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	var addCoupon model.Coupons

	err := utils.RequestDecoding(ctx, &addCoupon)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	err = validation.CheckValidation(&addCoupon)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}

	coupons.AddCouponService(ctx, &addCoupon)

}

// @Summary  		Get Coupon Handler
// @Description  	This Handler will get active coupons by passing "couponName" query parameters
// @Tags 			Coupon
// @Accept 			json
// @Procedure 		json
// @Param			couponName query string true "coupon name" SchemaExample({"couponName" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/coupon-get [get]
func GetCouponsHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	coupons.GetCouponsService(ctx)
}
