package coupons

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Add valid coupons to database
func AddCouponService(ctx *gin.Context, addCoupon *model.Coupons) {
	if db.RecordExist("coupons", "coupon_name", addCoupon.CouponName) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Coupon already exits")
		return
	}
	addCoupon.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	err := db.CreateRecord(&addCoupon)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "unable to create record")
		return
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Coupon added successfully",
		addCoupon,
		ctx,
	)
}

// Get coupon service with form value
func GetCouponsService(ctx *gin.Context) {
	couponName := ctx.Query("couponName")
	if couponName == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Coupon name is empty")
		return
	}

	if !db.RecordExist("coupons", "coupon_name", couponName) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Coupon Name does not exist")
		return
	}

	var coupons model.Coupons
	err := db.FindById(&coupons, couponName, "coupon_name")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Coupon not found")
		return
	}
	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Coupons detials are available",
		coupons,
		ctx,
	)
}
