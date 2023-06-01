package order

import (
	"fmt"
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/stripeservice"
	"main/server/services/token"
	"main/server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Get ID from token
func IdFromToken(ctx *gin.Context) (string, error) {
	tokenString, err := utils.GetTokenFromAuthHeader(ctx)
	if tokenString == "" {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Token not found")
		return "", err
	}
	if err != nil {

		return "", err
	}
	claims, err := token.DecodeToken(ctx, tokenString)
	if err != nil {

		return "", err
	}
	return claims.UserId, nil
}

// Make Payment for a specific product Service
func MakePaymentService(ctx *gin.Context, paymentRequest context.OrderRequest) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Error decoding token or invalid token")
		return
	}

	addressType := ctx.Query("addresstype")
	if addressType == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Address specified")
		return
	}

	var cartProduct model.CartProducts
	if !db.RecordExist("cart_products", "cart_id", paymentRequest.CartId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Cart id not found")
		return
	}

	if !db.RecordExist("cart_products", "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product id not found")
		return
	}
	err = db.FindById(&cartProduct, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving cart Details with given cart_id")
		return
	}
	if db.BothExists("cart_products", "cart_id", paymentRequest.CartId, "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "This cart don't have any such product")
		return
	}

	//coupons
	var productPrice float64
	var couponToRedeem model.Coupons
	if paymentRequest.CouponName != "" {
		productPrice, couponToRedeem, err = RedeemCoupon(ctx, paymentRequest.CouponName, cartProduct.ProductPrice, userId)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
			return
		}
	} else {
		productPrice = cartProduct.ProductPrice
	}

	// stripe payment api call
	pi, pi1 := stripeservice.StripePayment(int64(productPrice), paymentRequest.CardNumber, paymentRequest.ExpMonth, paymentRequest.ExpYear, paymentRequest.CVC, ctx)
	fmt.Println("pi", pi.Status)
	fmt.Println("pi1", pi1)

	//create payment
	var payment model.Payment
	payment.PaymentId = pi1.ID
	payment.UserId = userId
	payment.PaymentAmount = productPrice
	payment.PaymentType = "card"
	payment.PaymentStatus = string(pi1.Status)
	err = db.CreateRecord(&payment)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	//create order
	var order model.Order
	order.OrderId = payment.OrderId
	order.CartId = paymentRequest.CartId
	order.ProductId = paymentRequest.ProductId
	order.PaymentId = pi1.ID
	order.UserId = userId
	order.OrderQuantity = cartProduct.ProductCount
	order.OrderStatus = "CONFIRMED"
	order.OrderDate = time.Now().Format("2006-January-02")
	address, err := AlotAddressForConfirmedOrders(ctx, userId, addressType)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}
	order.OrderAddress = address
	err = db.CreateRecord(&order)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	//coupon redmeption update
	if paymentRequest.CouponName != "" {
		CouponRedemptionUpdate(couponToRedeem.CouponId, order.OrderId, order.CreatedAt)
	}

	//create user payment details
	var userPaymentDetails model.UserPayments
	userPaymentDetails.PaymentId = payment.PaymentId
	userPaymentDetails.UserId = userId
	userPaymentDetails.OrderId = payment.OrderId
	userPaymentDetails.CouponRedeemed = paymentRequest.CouponName
	err = db.CreateRecord(&userPaymentDetails)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	// make order request
	var orderRequest model.OrderRequest
	orderRequest.OrderId = order.OrderId
	orderRequest.UserId = userId
	orderRequest.OrderStatus = "DISPATCHED"
	if !db.RecordExist("db_constants", "constant_name", "dispatched") {
		var dbconstant model.DbConstant
		dbconstant.ConstantName = "dispatched"
		dbconstant.ConstantShortHand = "DISPATCHED"
		err = db.CreateRecord(&dbconstant)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "unable to create record")
			return
		}
	}
	err = db.CreateRecord(&orderRequest)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "unable to create record")
		return
	}

	//product details
	var productDetails model.Products
	err = db.FindById(&productDetails, paymentRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving cart Details with given cart_id")
		return
	}

	//inventory update
	var inventory model.Products
	err = db.FindById(&inventory, paymentRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error Finding record in DB ")
		return
	}
	inventory.ProductInventory--
	db.UpdateRecord(&inventory, paymentRequest.ProductId, "product_id")

	orderCompleteData := &response.OrderCompletionResponse{
		OrderId:         payment.OrderId,
		UserId:          userId,
		PaymentId:       pi1.ID,
		PaymentAmount:   productPrice,
		PaymentDate:     payment.CreatedAt,
		CartId:          paymentRequest.CartId,
		ProductId:       paymentRequest.ProductId,
		ProductName:     productDetails.ProductName,
		ProductCategory: productDetails.ProductCategory,
		ProductBrand:    productDetails.ProductBrand,
		OrderStatus:     order.OrderStatus,
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Congratulations your order has been created successfully",
		orderCompleteData,
		ctx,
	)

}

// function to redeem coupon with given code
func RedeemCoupon(ctx *gin.Context, couponName string, orderPrice float64, userId string) (float64, model.Coupons, error) {
	var couponToRedeem model.Coupons
	if !db.BothExists("user_payments", "user_id", userId, "coupon_redeemed", couponName) {
		return 0, couponToRedeem, fmt.Errorf("sorry , you already have redeemed this coupon")
	}
	if !db.RecordExist("coupons", "coupon_name", couponName) {
		return 0, couponToRedeem, fmt.Errorf("coupon does not exist")
	}

	err := db.FindById(&couponToRedeem, couponName, "coupon_name")
	if err != nil {
		return 0, couponToRedeem, fmt.Errorf("coupon not found from database")
	}

	productPrice := orderPrice * (1 - couponToRedeem.CouponPrice)

	return productPrice, couponToRedeem, nil
}

// Coupon Redemption table update
func CouponRedemptionUpdate(couponId string, orderId string, redeemedAt time.Time) {
	var couponRedemption model.CouponRedemptions
	couponRedemption.CouponId = couponId
	couponRedemption.OrderId = orderId
	couponRedemption.RedeemedAt = redeemedAt
	err := db.CreateRecord(&couponRedemption)
	if err != nil {
		return
	}
}

// Get Order Details that already has payment done
func GetOrderDetails(ctx *gin.Context) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Error decoding token or invalid token")
		return
	}

	var paymentInfo model.Payment
	err = db.FindById(&paymentInfo, userId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving user details")
		return
	}

	var orderInfo model.Order
	err = db.FindById(&orderInfo, userId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving user details")
		return
	}

	var productInfo model.Products
	err = db.FindById(&productInfo, orderInfo.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving user details")
		return
	}

	orderCompleteData := &response.OrderCompletionResponse{
		OrderId:         paymentInfo.OrderId,
		UserId:          userId,
		PaymentId:       paymentInfo.PaymentId,
		PaymentAmount:   paymentInfo.PaymentAmount,
		PaymentDate:     paymentInfo.CreatedAt,
		CartId:          orderInfo.CartId,
		OrderStatus:     orderInfo.OrderStatus,
		ProductId:       orderInfo.ProductId,
		ProductName:     productInfo.ProductName,
		ProductCategory: productInfo.ProductCategory,
		ProductBrand:    productInfo.ProductBrand,
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Here are your order details",
		orderCompleteData,
		ctx,
	)

}

// Alot address according to params passed
func AlotAddressForConfirmedOrders(ctx *gin.Context, userId string, addressType string) (string, error) {
	var userDefaultAddress model.UserAddresses
	query := "SELECT * FROM user_addresses WHERE user_id='" + userId + "' AND address_type='" + addressType + "'"
	err := db.QueryExecutor(query, &userDefaultAddress)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding user from DB")
		return "", err
	}

	address := userDefaultAddress.Name + " , " + userDefaultAddress.Street + " , " + userDefaultAddress.City + " , " + userDefaultAddress.State + " , " + userDefaultAddress.Country + " ,zip: " + userDefaultAddress.PostalCode + " , ph: " + userDefaultAddress.Phone

	return address, nil
}

// Cancel Order and Refund
func CancelOrderService(ctx *gin.Context, cancelOrderRequest context.CancelOrderRequest) {

	var order model.Order
	var payment model.Payment
	err := db.FindById(&order, cancelOrderRequest.OrderId, "order_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error fetching order details")
		return
	}

	if order.OrderStatus != "CONFIRMED" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Order not confirmed")
		return
	}

	err = db.FindById(&payment, order.PaymentId, "payment_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error fetching payment details")
		return
	}
	order.OrderStatus = "CANCELLED"
	query := "UPDATE orders SET order_status = 'CANCELLED' WHERE order_id ='" + order.OrderId + "'"
	err = db.QueryExecutor(query, &order)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Unable to execute query")
		return
	}

	// amount refund
	payment.RefundAmount = payment.PaymentAmount - payment.PaymentAmount*2.7
	query = "UPDATE payments set refund_amount = '" + strconv.Itoa(int(payment.RefundAmount)) + "' WHERE payment_id ='" + payment.PaymentId + "'"
	err = db.QueryExecutor(query, &payment)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error executing query")
		return
	}

	var userPayment model.UserPayments
	err = db.Delete(&userPayment, payment.PaymentId, "payment_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error deleting payment record")
		return
	}
}

// Make Cart Payment
func MakeCartPaymentService(ctx *gin.Context, paymentRequest context.CartOrderRequest) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Error decoding token or invalid token")
		return
	}

	addressType := ctx.Query("addresstype")
	if addressType == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Address specified")
		return
	}

	var cartProduct model.Cart
	if db.RecordExist("cart_products", "cart_id", paymentRequest.CartId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Order with the same cart id")
		return
	}

	err = db.FindById(&cartProduct, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving cart Details with given cart_id")
		return
	}

	//coupons
	var productPrice float64
	var couponToRedeem model.Coupons
	if paymentRequest.CouponName != "" {
		productPrice, couponToRedeem, err = RedeemCoupon(ctx, paymentRequest.CouponName, cartProduct.TotalPrice, userId)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
			return
		}
	} else {
		productPrice = cartProduct.TotalPrice
	}

	// stripe payment api call
	pi, pi1 := stripeservice.StripePayment(int64(productPrice), paymentRequest.CardNumber, paymentRequest.ExpMonth, paymentRequest.ExpYear, paymentRequest.CVC, ctx)
	fmt.Println("pi", pi.Status)
	fmt.Println("pi1", pi1)

	//create payment
	var payment model.Payment
	payment.PaymentId = pi1.ID
	payment.UserId = userId
	payment.PaymentAmount = productPrice
	payment.PaymentType = "card"
	payment.PaymentStatus = string(pi1.Status)
	err = db.CreateRecord(&payment)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	//create order
	var order model.Order
	order.OrderId = payment.OrderId
	order.CartId = paymentRequest.CartId
	order.PaymentId = pi1.ID
	order.UserId = userId
	order.OrderQuantity = 1
	order.OrderStatus = "CONFIRMED"
	order.OrderDate = time.Now().Format("utils.HTTP_OK 6-January-02")
	address, err := AlotAddressForConfirmedOrders(ctx, userId, addressType)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, err.Error())
		return
	}
	order.OrderAddress = address
	err = db.CreateRecord(&order)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	//coupon redmeption update
	if paymentRequest.CouponName != "" {
		CouponRedemptionUpdate(couponToRedeem.CouponId, order.OrderId, order.CreatedAt)
	}

	//create user payment details
	var userPaymentDetails model.UserPayments
	userPaymentDetails.PaymentId = payment.PaymentId
	userPaymentDetails.UserId = userId
	userPaymentDetails.OrderId = payment.OrderId
	userPaymentDetails.CouponRedeemed = paymentRequest.CouponName
	err = db.CreateRecord(&userPaymentDetails)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	// inventory update
	var inventory model.Products
	var cartProducts []model.CartProducts
	var productIds []string
	err = db.FindById(&cartProducts, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Unable to find cart id")
		return
	}
	for _, product := range cartProducts {
		productIds = append(productIds, product.ProductId)
	}
	for _, productId := range productIds {
		err = db.FindById(&inventory, productId, "product_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error finding inventory")
			return
		}
		inventory.ProductInventory--
		db.UpdateRecord(&inventory, productId, "product_id")
	}

	orderCompleteData := &response.CartOrderCompletionResponse{
		OrderId:       payment.OrderId,
		UserId:        userId,
		PaymentId:     pi1.ID,
		PaymentAmount: cartProduct.TotalPrice,
		PaymentDate:   payment.CreatedAt,
		CartId:        paymentRequest.CartId,
		OrderStatus:   order.OrderStatus,
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Congratulations your order has been created successfully",
		orderCompleteData,
		ctx,
	)

}

// Vendor Order Status Set
func VendorOrderStatusUpdateService(ctx *gin.Context, orderUpdateRequest context.VendorOrderStatusUpdate) {

	if !db.RecordExist("order_request", "order_id", orderUpdateRequest.OrderId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Order does not exist")
		return
	}
	var orderRequest model.OrderRequest

	err := db.FindById(&orderRequest, orderUpdateRequest.OrderId, "order_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Unable to find order request")
		return
	}
	orderRequest.OrderStatus = "DELIVERED"
	db.UpdateRecord(&orderRequest, orderUpdateRequest.OrderId, "order_id")

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Order Status updated",
		orderRequest,
		ctx,
	)

}
