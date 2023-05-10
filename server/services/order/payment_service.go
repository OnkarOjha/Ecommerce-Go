package order

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"
	"strconv"
	"time"
)

func UserIdFromToken(ctx *gin.Context) (string, error) {
	tokenString, err := utils.GetTokenFromAuthHeader(ctx)
	if err != nil {

		return "", err
	}
	claims, err := token.DecodeToken(ctx, tokenString)
	if err != nil {

		return "", err
	}
	return claims.UserId, nil
}

func StripePayment(OrderPrice int64, cardNumber string, expMonth string, expYear string, cvc string, ctx *gin.Context) (pi, pi1 *stripe.PaymentIntent) {
	//TODO
	stripe.Key = "sk_test_51MvCYxSGxKXiPagaKdfa8MM2nYhjysJ41IUESqCLjca0meSTlzal4wbqMFZDbpTa5w1YXvdwygU8yMYbBecfgLCC00Yrx2WfFF"

	pm, err := paymentmethod.New(&stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cardNumber),
			ExpMonth: stripe.String(expMonth),
			ExpYear:  stripe.String(expYear),
			CVC:      stripe.String(cvc),
		},
	})
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error creating card details")
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(OrderPrice) * 100),
		Currency:           stripe.String("inr"),
		Description:        stripe.String("Payment"),
		PaymentMethod:      stripe.String(pm.ID),
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
	}
	pi, err = paymentintent.New(params)
	if err != nil {
		response.ErrorResponse(ctx, 503, "Error processing payment")
		return
	}

	params1 := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String("pm_card_visa"),
		CaptureMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
	}

	pi1, err = paymentintent.Confirm(pi.ID, params1)
	if err != nil {
		response.ErrorResponse(ctx, 503, "Error confirming payment")
		return
	}

	switch pi1.Status {
	case "succeeded":
		response.ShowResponse("Success", utils.HTTP_OK, "Payment processed Successfully", "", ctx)
		return
	case "requires_payment_method":
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Requires Payment Method")
		return
	case "requires_action":

		if pi1.Status == "requires_action" && pi1.NextAction != nil {
			switch pi1.NextAction.Type {
			case "use_stripe_sdk":

				response.ShowResponse(
					"Success",
					utils.HTTP_OK,
					"Payment processed Successfully , Here is your client secret",
					pi1.ClientSecret,
					ctx,
				)
			}
		}
	default:
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Payment requires more actions")
		return
	}

	return pi, pi1

}

func MakePaymentService(ctx *gin.Context, paymentRequest context.OrderRequest) {
	userId, err := UserIdFromToken(ctx)
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
	pi, pi1 := StripePayment(int64(cartProduct.ProductPrice), paymentRequest.CardNumber, paymentRequest.ExpMonth, paymentRequest.ExpYear, paymentRequest.CVC, ctx)
	fmt.Println("pi", pi.Status)
	fmt.Println("pi1", pi1)

	var payment model.Payment
	//create payment
	payment.PaymentId = pi1.ID
	payment.UserId = userId
	payment.PaymentAmount = cartProduct.ProductPrice
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

	//create user payment details
	var userPaymentDetails model.UserPayments
	userPaymentDetails.PaymentId = payment.PaymentId
	userPaymentDetails.UserId = userId
	userPaymentDetails.OrderId = payment.OrderId
	err = db.CreateRecord(&userPaymentDetails)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	//product details
	var productDetails model.Products
	err = db.FindById(&productDetails, paymentRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving cart Details with given cart_id")
		return
	}

	orderCompleteData := &response.OrderCompletionResponse{
		OrderId:         payment.OrderId,
		UserId:          userId,
		PaymentId:       pi1.ID,
		PaymentAmount:   cartProduct.ProductPrice,
		PaymentDate:     payment.CreatedAt,
		CartId:          paymentRequest.CartId,
		ProductId:       paymentRequest.ProductId,
		ProductName:     productDetails.ProductName,
		ProductCategory: productDetails.ProductCategory,
		ProductBrand:    productDetails.ProductBrand,
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Congratulations your order has been created successfully",
		orderCompleteData,
		ctx,
	)

}

func GetOrderDetails(ctx *gin.Context) {
	userId, err := UserIdFromToken(ctx)
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
	db.QueryExecutor(query, &order)

	// amount refund
	payment.RefundAmount = payment.PaymentAmount - payment.PaymentAmount*2.7
	query = "UPDATE payments set refund_amount = '" + strconv.Itoa(int(payment.RefundAmount)) + "' WHERE payment_id ='" + payment.PaymentId + "'"
	db.QueryExecutor(query, &payment)

	var userPayment model.UserPayments
	db.Delete(&userPayment, payment.PaymentId, "payment_id")
}

func MakeCartPaymentService(ctx *gin.Context, paymentRequest context.CartOrderRequest) {
	userId, err := UserIdFromToken(ctx)
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
	// if db.RecordExist("cart_products", "cart_id", paymentRequest.CartId) {
	// 	response.ErrorResponse(context, utils.HTTP_BAD_REQUEST , "Order with the same cart id")
	// 	return
	// }
	fmt.Println("fkjfj", paymentRequest)
	err = db.FindById(&cartProduct, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error retrieving cart Details with given cart_id")
		return
	}

	pi, pi1 := StripePayment(int64(cartProduct.TotalPrice), paymentRequest.CardNumber, paymentRequest.ExpMonth, paymentRequest.ExpYear, paymentRequest.CVC, ctx)
	fmt.Println("pi", pi.Status)
	fmt.Println("pi1", pi1)

	var payment model.Payment
	//create payment
	payment.PaymentId = pi1.ID
	payment.UserId = userId
	payment.PaymentAmount = cartProduct.TotalPrice
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

	//create user payment details
	var userPaymentDetails model.UserPayments
	userPaymentDetails.PaymentId = payment.PaymentId
	userPaymentDetails.UserId = userId
	userPaymentDetails.OrderId = payment.OrderId
	err = db.CreateRecord(&userPaymentDetails)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error creating record: "+err.Error())
		return
	}

	orderCompleteData := &response.CartOrderCompletionResponse{
		OrderId:       payment.OrderId,
		UserId:        userId,
		PaymentId:     pi1.ID,
		PaymentAmount: cartProduct.TotalPrice,
		PaymentDate:   payment.CreatedAt,
		CartId:        paymentRequest.CartId,
	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Congratulations your order has been created successfully",
		orderCompleteData,
		ctx,
	)

}
