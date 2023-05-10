package handler

import (
	"main/server/model"
	"main/server/response"
	"main/server/services/coupons"
	"main/server/utils"
	"main/server/validation"

	"github.com/gin-gonic/gin"
)

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

func GetCouponsHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	coupons.GetCouponsService(ctx)
}
