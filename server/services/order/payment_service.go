package order

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/provider"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
)

func UserIdFromToken(context *gin.Context) (string, error) {
	tokenString, err := utils.GetTokenFromAuthHeader(context)
	if err != nil {

		return "", err
	}
	claims, err := provider.DecodeToken(context, tokenString)
	if err != nil {

		return "", err
	}
	return claims.UserId, nil
}

func StripePayment(OrderPrice int64, cardNumber string, expMonth string, expYear string, cvc string, context *gin.Context) (pi, pi1 *stripe.PaymentIntent) {
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
		response.ErrorResponse(context, 400, "Error creating card details")
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
		response.ErrorResponse(context, 503, "Error processing payment")
		return
	}

	params1 := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String("pm_card_visa"),
		CaptureMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
	}

	pi1, err = paymentintent.Confirm(pi.ID, params1)
	if err != nil {
		response.ErrorResponse(context, 503, "Error confirming payment")
		return
	}

	switch pi1.Status {
	case "succeeded":
		response.ShowResponse("Success", 200, "Payment processed Successfully", "", context)
		return
	case "requires_payment_method":
		response.ErrorResponse(context, 400, "Requires Payment Method")
		return
	case "requires_action":

		if pi1.Status == "requires_action" && pi1.NextAction != nil {
			switch pi1.NextAction.Type {
			case "use_stripe_sdk":

				response.ShowResponse(
					"Success",
					200,
					"Payment processed Successfully , Here is your client secret",
					pi1.ClientSecret,
					context,
				)
			}
		}
	default:
		response.ErrorResponse(context, 400, "Payment requires more actions")
		return
	}

	return pi, pi1

}

func MakePaymentService(context *gin.Context, paymentRequest request.OrderRequest) {
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 401, "Error decoding token or invalid token")
		return
	}

	var cartProduct model.CartProducts
	if !db.RecordExist("cart_products", "cart_id", paymentRequest.CartId) {
		response.ErrorResponse(context, 400, "Cart id not found")
		return
	}

	if !db.RecordExist("cart_products", "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(context, 400, "Product id not found")
		return
	}
	err = db.FindById(&cartProduct, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving cart Details with given cart_id")
		return
	}
	if db.BothExists("cart_products", "cart_id", paymentRequest.CartId, "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(context, 400, "This cart don't have any such product")
		return
	}
	pi, pi1 := StripePayment(int64(cartProduct.ProductPrice), paymentRequest.CardNumber, paymentRequest.ExpMonth, paymentRequest.ExpYear, paymentRequest.CVC, context)
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
		response.ErrorResponse(context, 500, "Error creating record: "+err.Error())
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
	order.OrderStatus = "Booked"
	order.OrderDate = time.Now().Format("2006-January-02")
	err = db.CreateRecord(&order)
	if err != nil {
		response.ErrorResponse(context, 500, "Error creating record: "+err.Error())
		return
	}

	//create user payment details
	var userPaymentDetails model.UserPayments
	userPaymentDetails.PaymentId = payment.PaymentId
	userPaymentDetails.UserId = userId
	userPaymentDetails.OrderId = payment.OrderId
	err = db.CreateRecord(&userPaymentDetails)
	if err != nil {
		response.ErrorResponse(context, 500, "Error creating record: "+err.Error())
		return
	}

	//product details
	var productDetails model.Products
	err = db.FindById(&productDetails, paymentRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving cart Details with given cart_id")
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
		200,
		"Congratulations your order has been created successfully",
		orderCompleteData,
		context,
	)

}

func GetOrderDetails(context *gin.Context) {
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 401, "Error decoding token or invalid token")
		return
	}

	var paymentInfo model.Payment
	err = db.FindById(&paymentInfo, userId, "user_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving user details")
		return
	}

	var orderInfo model.Order
	err = db.FindById(&orderInfo, userId, "user_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving user details")
		return
	}

	var productInfo model.Products
	err = db.FindById(&productInfo, orderInfo.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving user details")
		return
	}

	orderCompleteData := &response.OrderCompletionResponse{
		OrderId:         paymentInfo.OrderId,
		UserId:          userId,
		PaymentId:       paymentInfo.PaymentId,
		PaymentAmount:   paymentInfo.PaymentAmount,
		PaymentDate:     paymentInfo.CreatedAt,
		CartId:          orderInfo.CartId,
		ProductId:       orderInfo.ProductId,
		ProductName:     productInfo.ProductName,
		ProductCategory: productInfo.ProductCategory,
		ProductBrand:    productInfo.ProductBrand,
	}

	response.ShowResponse(
		"Success",
		200,
		"Here are your order details",
		orderCompleteData,
		context,
	)

}
