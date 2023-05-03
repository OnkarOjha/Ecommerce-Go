package order

import (



	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"

	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"

	// "github.com/gin-gonic/gin"
	// "github.com/stripe/stripe-go/v72"
	// "github.com/stripe/stripe-go/v72/paymentintent"

	// // "github.com/stripe/stripe-go/v72/card"
	// "github.com/stripe/stripe-go/v72/paymentmethod"
)
func StripePayment(OrderPrice int64, context *gin.Context)  {
	//TODO
	stripe.Key = "sk_test_51MvCYxSGxKXiPagaKdfa8MM2nYhjysJ41IUESqCLjca0meSTlzal4wbqMFZDbpTa5w1YXvdwygU8yMYbBecfgLCC00Yrx2WfFF"
	pi, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(int64(OrderPrice) * 100),
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		
	})
	fmt.Println("pi: ",pi)
	if err != nil {
		// c.String(http.StatusBadRequest, "Payment Unsuccessfull")
		return
	}

}
// func StripePayment(OrderPrice int64, context *gin.Context) (pi, pi1 *stripe.PaymentIntent) {
// 	//TODO
// 	stripe.Key = "sk_test_51MvCYxSGxKXiPagaKdfa8MM2nYhjysJ41IUESqCLjca0meSTlzal4wbqMFZDbpTa5w1YXvdwygU8yMYbBecfgLCC00Yrx2WfFF"

// 	pm, err := paymentmethod.New(&stripe.PaymentMethodParams{ 	
// 		Type: stripe.String("card"),
// 		Card: &stripe.PaymentMethodCardParams{
// 			Number:   stripe.String("4242424242424242"),
// 			ExpMonth: stripe.String("12"),
// 			ExpYear:  stripe.String("2024"),
// 			CVC:      stripe.String("123"),
// 		},
// 	})
// 	if err!=nil{
// 		response.ErrorResponse(context , 400 , "Error creating card details")
// 		return
// 	}

// 	params := &stripe.PaymentIntentParams{
// 		Amount:        stripe.Int64(int64(OrderPrice) * 100),
// 		Currency:      stripe.String("inr"),
// 		Description:   stripe.String("Payment"),
// 		PaymentMethod: stripe.String(pm.ID),
// 		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
// 	}
// 	pi, err = paymentintent.New(params)
// 	if err != nil {
// 		response.ErrorResponse(context, 503, "Error processing payment")
// 		return
// 	}

// 	params1 := &stripe.PaymentIntentConfirmParams{
// 		PaymentMethod: stripe.String("pm_card_visa"),
// 		CaptureMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodAutomatic)),
// 	}
	

// 	pi1, err = paymentintent.Confirm(pi.ID, params1)
// 	if err != nil {
// 		response.ErrorResponse(context, 503, "Error confirming payment")
// 		return
// 	}

// 	// Check the payment intent status
// 	switch pi1.Status {
// 	case "succeeded":
// 		response.ShowResponse("Success", 200, "Payment processed Successfully", "", context)
// 		return
// 	case "requires_payment_method":
// 		response.ErrorResponse(context, 400, "Requires Payment Method")
// 		return
// 	case "requires_action":
		
// 		if pi1.Status == "requires_action" && pi1.NextAction != nil {
// 			switch pi1.NextAction.Type {
// 			case "use_stripe_sdk":
// 				fmt.Println("oo ethe aa")
				
// 			}
// 		}
// 		// // Additional action required
// 		// clientSecret := pi1.ClientSecret
// 		// w.Header().Set("Content-Type", "application/json")
// 		// w.WriteHeader(http.StatusOK)
// 		// json.NewEncoder(w).Encode(map[string]string{
// 		// 	"client_secret": clientSecret,
// 		// })
// 	default:
// 		// Unknown status
// 		response.ErrorResponse(context, 400, "Payment requires more actions")
// 		return
// 	}

// 	return pi, pi

// }

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
	// pi, pi1 := StripePayment(int64(cartProduct.ProductPrice), context)
	// fmt.Println("pi", pi)
	// fmt.Println("pi1", pi1)
	StripePayment(int64(cartProduct.ProductPrice), context)
}

// func ConfirmPayment(context *gin.Context){
	
// 	// Create a PaymentIntentParams struct with the client secret
// 	params := &stripe.PaymentIntentParams{
// 		ClientSecret: stripe.String("pi_3N3Gk9SGxKXiPaga10SPiiu2_secret_9WIrXYAWLHngUpRdFyYBtH4WI"),
// 	}

// 	// Confirm the payment intent with the params
// 	pi, err := paymentintent.Confirm(params)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the ID of the confirmed payment intent
// 	fmt.Println("Confirmed payment intent ID:", pi.ID)
// }
