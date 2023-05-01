package cart

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/provider"
	"main/server/request"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func UserIdFromToken(context *gin.Context) string {
	tokenString, err := utils.GetTokenFromAuthHeader(context)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		context.Abort()
	}
	claims, err := provider.DecodeToken(tokenString)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		context.Abort()
	}
	return claims.UserId
}

func AddToCartService(context *gin.Context, addToCartRequest request.AddToCartRequest) {
	if !db.RecordExist("products", "product_id", addToCartRequest.ProductId) {
		response.ErrorResponse(context, 400, "Invalid Product ID")
		return
	}
	var product model.Products
	var cartProduct model.CartProducts
	var cart model.Cart

	if db.RecordExist("cart_products", "product_id", addToCartRequest.ProductId) {
		response.ErrorResponse(context, 400, "Product already added to cart if you want to add more please proceed to /addProduct")
		return
	}
	cartProduct.ProductId = addToCartRequest.ProductId
	cartProduct.ProductCount = addToCartRequest.ProductCount

	// fetch product price from products table
	err := db.FindById(&product, addToCartRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Product not found")
		return
	}
	var cartPreviousProduct model.CartProducts
	if !db.RecordExist("cart_products", "user_id", addToCartRequest.UserId) {
		err = db.FindById(&cartPreviousProduct, addToCartRequest.UserId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Product not found")
			return
		}
		cartProduct.CartId = cartPreviousProduct.CartId
	}
	cartProduct.ProductPrice = product.ProductPrice * addToCartRequest.ProductCount
	cartProduct.UserId = addToCartRequest.UserId
	err = db.CreateRecord(&cartProduct)
	if err != nil {
		response.ErrorResponse(context, 500, err.Error())
		return
	}

	// cart table
	cart.CartId = cartProduct.CartId
	cart.UserId = addToCartRequest.UserId
	cart.CartCount = cart.CartCount + 1
	cart.TotalPrice = cart.TotalPrice + cartProduct.ProductPrice

	if !db.RecordExist("carts", "user_id", addToCartRequest.UserId) {
		fmt.Println("New cart creation")
		err = db.CreateRecord(&cart)
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}
	} else {
		var cartUpdate model.Cart
		err := db.FindById(&cartUpdate, addToCartRequest.UserId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Record not found")
			return
		}
		cart.CartId = cartProduct.CartId
		cart.UserId = addToCartRequest.UserId
		cart.CartCount = cartUpdate.CartCount + 1
		cart.TotalPrice = cartUpdate.TotalPrice + cartProduct.ProductPrice
		db.UpdateRecord(&cart, addToCartRequest.UserId, "user_id")
	}
	response.ShowResponse(
		"Success",
		200,
		"Product added to cart",
		cartProduct,
		context,
	)
	response.ShowResponse(
		"Success",
		200,
		"Cart details updated successfully",
		cart,
		context,
	)
}

func AddProductService(context *gin.Context, addProductCountRequest request.AddToCartRequest) {
	var cartProduct model.CartProducts
	var product model.Products
	var cart model.Cart
	if db.RecordExist("cart_products", "product_id", addProductCountRequest.ProductId) {

		err := db.FindById(&cartProduct, addProductCountRequest.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Record not found")
			return
		}
		err = db.FindById(&product, addProductCountRequest.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Product not found")
			return
		}
		cartProduct.ProductCount = cartProduct.ProductCount + addProductCountRequest.ProductCount

		cartProduct.ProductPrice += addProductCountRequest.ProductCount * product.ProductPrice

		err = db.UpdateRecord(&cartProduct, addProductCountRequest.ProductId, "product_id").Error
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}

		err = db.FindById(&cart, addProductCountRequest.UserId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}

		cart.TotalPrice = cart.TotalPrice + cartProduct.ProductPrice

		err = db.UpdateRecord(&cart, cart.CartId, "cart_id").Error
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}

		response.ShowResponse(
			"Success",
			200,
			"Product added successfully",
			cartProduct,
			context,
		)
		response.ShowResponse(
			"Success",
			200,
			"Cart Updated successfully",
			cart,
			context,
		)
	}
}

func RemoveFromCartService(context *gin.Context, removeFromCartRequest request.RemoveFromCart) {
	if !db.RecordExist("cart_products", "cart_id", removeFromCartRequest.CartId) {
		response.ErrorResponse(context, 400, "Cart Id not found")
		return
	}
	if !db.RecordExist("cart_products", "product_id", removeFromCartRequest.ProductId) {
		response.ErrorResponse(context, 400, "Product Id not found")
		return
	}
	var cartProduct model.CartProducts
	var cart model.Cart
	err := db.FindById(&cartProduct, removeFromCartRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error retrieving cart details from cart_products")
		return
	}

	err = db.FindById(&cart, removeFromCartRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error retrieving cart details from cart")
		return
	}

	cart.TotalPrice = cart.TotalPrice - cartProduct.ProductPrice
	cart.CartCount = cart.CartCount - 1

	err = db.UpdateRecord(&cart, removeFromCartRequest.CartId, "cart_id").Error
	if err != nil {
		response.ErrorResponse(context, 500, "Error updating cart")
		return
	}

	err = db.DeleteRecord(&cartProduct, removeFromCartRequest.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error Deleting product from cart")
		return
	}
	response.ShowResponse("Success", 200, "Product deleted", cartProduct, context)
	response.ShowResponse("Success", 200, "Cart Detials ", cart, context)
}

func RemoveProductService(context *gin.Context, removeProductFromCart request.RemoveProduct) {
	if !db.RecordExist("cart_products", "cart_id", removeProductFromCart.CartId) {
		response.ErrorResponse(context, 400, "Cart Id not found")
		return
	}
	if !db.RecordExist("cart_products", "product_id", removeProductFromCart.ProductId) {
		response.ErrorResponse(context, 400, "Product Id not found")
		return
	}

	var product model.Products
	var cartProduct model.CartProducts
	var cart model.Cart
	err := db.FindById(&product, removeProductFromCart.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Product Id not found")
		return
	}

	err = db.FindById(&cartProduct, removeProductFromCart.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Cart Id not found")
		return
	}
	if cartProduct.ProductCount < removeProductFromCart.ProductCount {
		response.ErrorResponse(context, 400, "Request limit exceeded")
		return
	}

	cartProduct.ProductCount -= removeProductFromCart.ProductCount
	cartProduct.ProductPrice -= product.ProductPrice * removeProductFromCart.ProductCount

	if cartProduct.ProductCount <= 0 {

		cartProduct.ProductPrice = 0.0
	}

	query := "UPDATE cart_products SET product_price= 0 , product_count = 0  where cart_id=?"
	db.QueryExecutor(query, &cartProduct, removeProductFromCart.CartId)

	err = db.FindById(&cart, removeProductFromCart.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error finding cart details")
		return
	}
	cart.TotalPrice -= product.ProductPrice * removeProductFromCart.ProductCount
	if cart.TotalPrice <= 0 {
		cart.TotalPrice = 0.0
	}

	query = "UPDATE carts SET total_price = 0 where cart_id = ?"
	db.QueryExecutor(query, &cart, removeProductFromCart.CartId)

	if cartProduct.ProductCount == 0 {
		err = db.DeleteRecord(&cartProduct, removeProductFromCart.CartId, "cart_id")
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}
	}

	response.ShowResponse("Success", 200, "Successfully decremented the product count", cartProduct, context)
	response.ShowResponse("Success", 200, "Cart Details after decrement", cart, context)
}
