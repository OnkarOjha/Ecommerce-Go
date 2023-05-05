package cart

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/provider"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	// "main/server/response/cart_response"
	"github.com/gin-gonic/gin"
)

func UserIdFromToken(context *gin.Context) (string, error) {
	tokenString, err := utils.GetTokenFromAuthHeader(context)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		return "", err
	}
	claims, err := provider.DecodeToken(context, tokenString)
	if err != nil {
		response.ErrorResponse(
			context, 401, "Error decoding token or invalid token",
		)
		return "", err
	}
	return claims.UserId, nil
}

func AddToCartService(context *gin.Context, addToCartRequest request.AddToCartRequest) {
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 400, "Error in Token header , no userId found")
		return
	}
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
	err = db.FindById(&product, addToCartRequest.ProductId, "product_id")
	if err != nil {

		response.ErrorResponse(context, 400, "Product not found")
		return
	}
	var cartPreviousProduct model.CartProducts
	if db.RecordExist("cart_products", "user_id", userId) {
		err = db.FindById(&cartPreviousProduct, userId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Product not found")
			return
		}
		cartProduct.CartId = cartPreviousProduct.CartId
	}
	cartProduct.ProductPrice = product.ProductPrice * addToCartRequest.ProductCount
	cartProduct.UserId = userId
	err = db.CreateRecord(&cartProduct)
	if err != nil {
		response.ErrorResponse(context, 500, err.Error())
		return
	}

	// cart table
	cart.CartId = cartProduct.CartId
	cart.UserId = userId
	cart.CartCount = cart.CartCount + 1
	cart.TotalPrice = cart.TotalPrice + cartProduct.ProductPrice

	if !db.RecordExist("carts", "user_id", userId) {
		fmt.Println("New cart creation")
		err = db.CreateRecord(&cart)
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}
	} else {
		var cartUpdate model.Cart
		err := db.FindById(&cartUpdate, userId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 400, "Record not found")
			return
		}
		cart.CartId = cartProduct.CartId
		cart.UserId = userId
		cart.CartCount = cartUpdate.CartCount + 1
		cart.TotalPrice = cartUpdate.TotalPrice + cartProduct.ProductPrice
		db.UpdateRecord(&cart, userId, "user_id")
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
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 400, "Error in Token header , no userId found")
		return
	}
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

		err = db.FindById(&cart, userId, "user_id")
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}

		cart.TotalPrice += cartProduct.ProductPrice
		cart.TotalPrice -= product.ProductPrice
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
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 400, "Error in Token header , no userId found")
		return
	}
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
	err = db.FindById(&cartProduct, removeFromCartRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error retrieving product from cart_products")
		return
	}

	err = db.FindById(&cart, userId, "user_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error retrieving cart details from cart")
		return
	}

	cart.TotalPrice = cart.TotalPrice - cartProduct.ProductPrice
	cart.CartCount = cart.CartCount - 1

	query := "update carts set total_price = 0 , cart_count = 0 where cart_id=?"
	err = db.QueryExecutor(query, &cart, removeFromCartRequest.CartId)
	if err != nil {
		response.ErrorResponse(context, 400, "Error updating cart")
		return
	}

	err = db.Delete(&cartProduct, removeFromCartRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(context, 500, "Error Deleting product from cart")
		return
	}
	response.ShowResponse("Success", 200, "Product deleted", cartProduct, context)
	response.ShowResponse("Success", 200, "Cart Details", cart, context)
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

	err = db.UpdateRecord(&cartProduct, removeProductFromCart.CartId, "cart_id").Error
	if err != nil {
		response.ErrorResponse(context, 500, "Error updating cart")
	}

	err = db.FindById(&cart, removeProductFromCart.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(context, 400, "Error finding cart details")
		return
	}
	cart.TotalPrice -= product.ProductPrice * removeProductFromCart.ProductCount
	if cart.TotalPrice <= 0 {
		cart.TotalPrice = 0.0
		cart.CartCount = 0
	}
	cart.CartCount -= 1
	err = db.UpdateRecord(&cart, cart.CartId, "cart_id").Error
	if err != nil {
		response.ErrorResponse(context, 500, "Error updating cart")
	}

	if cartProduct.ProductCount == 0 {
		err = db.Delete(&cartProduct, removeProductFromCart.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(context, 500, err.Error())
			return
		}
	}

	response.ShowResponse("Success", 200, "Successfully decremented the product count", cartProduct, context)
	response.ShowResponse("Success", 200, "Cart Details after decrement", cart, context)
}

func GetCartDetailsService(context *gin.Context) {
	userId, err := UserIdFromToken(context)
	if err != nil {
		response.ErrorResponse(context, 400, "Error in Token header , no userId found")
		return
	}

	var cartProductDetails model.CartProducts
	err = db.FindById(&cartProductDetails, userId, "user_id")
	if err != nil {
		response.ErrorResponse(context, 401, "No cart products matching this user id")
		return
	}

	var cartResponse response.CartProductResponse

	cartResponse.CartId = cartProductDetails.CartId
	cartResponse.ProductId = cartProductDetails.ProductId
	cartResponse.ProductCount = cartProductDetails.ProductCount
	cartResponse.ProductPrice = cartProductDetails.ProductPrice
	cartResponse.ProductAddedAt = cartProductDetails.CreatedAt

	response.ShowResponse(
		"Success",
		200,
		"User cart details are shown below",
		cartResponse,
		context,
	)
}
