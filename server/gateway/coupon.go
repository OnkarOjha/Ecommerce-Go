package gateway

import (
	"main/server/db"
	"main/server/model"
	"time"

	"github.com/gin-gonic/gin"
)

//Coupon Session management middleware
func CouponExpiration(ctx *gin.Context) {
	var coupons []model.Coupons

	db.InitDB().Find(&coupons)

	for _, coupon := range coupons {
		if coupon.ExpiresAt.Before(time.Now()) {
			db.DeleteRecord(&coupon, coupon.CouponId, "coupon_id")
		}

	}

}
