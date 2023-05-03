package order

import (

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"

)

func UserIdFromToken(context *gin.Context) (string , error) {
	tokenString, err := utils.GetTokenFromAuthHeader(context)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		return "" , err
	}
	claims, err := provider.DecodeToken(tokenString)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		return "" , err
	}
	return claims.UserId , nil
}

func StripePayment(OrderPrice int64, context *gin.Context) (pi, pi1 *stripe.PaymentIntent) {
	//TODO
	stripe.Key = "sk_test_51MvCYxSGxKXiPagaKdfa8MM2nYhjysJ41IUESqCLjca0meSTlzal4wbqMFZDbpTa5w1YXvdwygU8yMYbBecfgLCC00Yrx2WfFF"

	pm, err := paymentmethod.New(&stripe.PaymentMethodParams{ 	
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.String("12"),
			ExpYear:  stripe.String("2024"),
			CVC:      stripe.String("123"),
		},
	})
	if err!=nil{
		response.ErrorResponse(context , 400 , "Error creating card details")
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(OrderPrice) * 100),
		Currency:      stripe.String("inr"),
		Description:   stripe.String("Payment"),
		PaymentMethod: stripe.String(pm.ID),
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
	//TODO
	var cartProduct model.CartProducts
	if !db.RecordExist("cart_products", "cart_id", paymentRequest.CartId) {
		response.ErrorResponse(context, 400, "Cart id not found")
		return
	}

	if !db.RecordExist("cart_products", "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(context, 400, "Product id not found")
		return
	}
	err := db.FindById(&cartProduct, paymentRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving cart Details with given cart_id")
		return
	}
	if db.BothExists("cart_products", "cart_id", paymentRequest.CartId, "product_id", paymentRequest.ProductId) {
		response.ErrorResponse(context, 400, "This cart don't have any such product")
		return
	}
	pi, pi1 := StripePayment(int64(cartProduct.ProductPrice), context)
	fmt.Println("pi", pi.Status)
	fmt.Println("pi1", pi1)

	var payment model.Payment

	
	userId , err := UserIdFromToken(context)
	if err!=nil{
		response.ErrorResponse(context, 401, "Error decoding token or invalid token")
        return
	}

	//create payment
	payment.PaymentId  = pi1.paymentId
	payment.UserId = userId
	payment.PaymentAmount = pi1.Amount/100
	payment.PaymentType = pi1.PaymentType
	payment.PaymentStatus = pi1.Status
	err = db.CreateRecord(&payment)
	if err != nil {
		response.ErrorResponse(context, 500, "Error creating record: "+err.Error())
		return1
	}


	//create order
	var order model.order
	order.OrderId = payment.OrderId
	order.CartId = paymentRequest.CartId
	order.ProductId = paymentRequest.ProductId
	order.UserId = userId
	order.OrderQuantity = cartProduct.ProductCount
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
	err := db.FindById(&productDetails, paymentRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error retrieving cart Details with given cart_id")
		return
	}

	response.CartProductResponse{
		"Success",
		200,
		"Congratulations your order has been created successfully",
		payment.OrderId,
		userId,
		pi1.paymentId,
		pi.Amount/100,
		payment.CreatedAt,
		paymentRequest.CartId,
		paymentRequest.ProductId,
        productDetails.ProductName,
		productDetails.ProductCategory,
		productDetails.ProductBrand,
	}

}

