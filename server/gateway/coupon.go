package gateway

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Coupon Session management middleware
func CouponExpiration(ctx *gin.Context) {
	var coupons []model.Coupons

	db.InitDB().Find(&coupons)

	for _, coupon := range coupons {
		if coupon.ExpiresAt.Before(time.Now()) {
			err := db.DeleteRecord(&coupon, coupon.CouponId, "coupon_id")
			if err != nil {
				response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Not able to delete DB record")
				return
			}
		}

	}

}
