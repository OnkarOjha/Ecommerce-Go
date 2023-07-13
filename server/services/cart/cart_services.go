package cart

import (
	"fmt"
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// Get ID from token
func IdFromToken(ctx *gin.Context) (string, error) {
	tokenString, err := utils.GetTokenFromAuthHeader(ctx)
	if err != nil {
		response.ErrorResponse(
			ctx, utils.HTTP_UNAUTHORIZED, "Error decoding token or invalid token",
		)
		return "", err
	}
	claims, err := token.DecodeToken(ctx, tokenString)
	if err != nil {
		response.ErrorResponse(
			ctx, utils.HTTP_UNAUTHORIZED, "Error decoding token or invalid token",
		)
		return "", err
	}
	return claims.UserId, nil
}

// Service to add Product to cart
func AddToCartService(ctx *gin.Context, addToCartRequest context.AddToCartRequest) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error in Token header , no userId found")
		return
	}
	if !db.RecordExist("products", "product_id", addToCartRequest.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Invalid Product ID")
		return
	}
	var product model.Products
	var cartProduct model.CartProducts
	var cart model.Cart

	if db.RecordExist("cart_products", "product_id", addToCartRequest.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product already added to cart if you want to add more please proceed to /addProduct")
		return
	}
	cartProduct.ProductId = addToCartRequest.ProductId
	cartProduct.ProductCount = addToCartRequest.ProductCount

	// fetch product price from products table
	err = db.FindById(&product, addToCartRequest.ProductId, "product_id")
	if err != nil {

		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product not found")
		return
	}
	var cartPreviousProduct model.CartProducts
	if db.RecordExist("cart_products", "user_id", userId) {
		err = db.FindById(&cartPreviousProduct, userId, "user_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product not found")
			return
		}
		cartProduct.CartId = cartPreviousProduct.CartId
	}
	cartProduct.ProductPrice = product.ProductPrice * addToCartRequest.ProductCount
	cartProduct.UserId = userId
	err = db.CreateRecord(&cartProduct)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	// update cart table
	cart.CartId = cartProduct.CartId
	cart.UserId = userId
	cart.CartCount = cart.CartCount + 1
	cart.TotalPrice = cart.TotalPrice + cartProduct.ProductPrice

	if !db.RecordExist("carts", "user_id", userId) {
		fmt.Println("New cart creation")
		err = db.CreateRecord(&cart)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
			return
		}
	} else {
		var cartUpdate model.Cart
		err := db.FindById(&cartUpdate, userId, "user_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Record not found")
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
		utils.HTTP_OK,
		"Product added to cart",
		cartProduct,
		ctx,
	)
	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"Cart details updated successfully",
		cart,
		ctx,
	)
}

// Add more products of same type to cart
func AddProductService(ctx *gin.Context, addProductCountRequest context.AddToCartRequest) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error in Token header , no userId found")
		return
	}
	var cartProduct model.CartProducts
	var product model.Products
	var cart model.Cart
	if db.RecordExist("cart_products", "product_id", addProductCountRequest.ProductId) {

		err := db.FindById(&cartProduct, addProductCountRequest.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Record not found")
			return
		}
		err = db.FindById(&product, addProductCountRequest.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product not found")
			return
		}
		cartProduct.ProductCount = cartProduct.ProductCount + addProductCountRequest.ProductCount

		cartProduct.ProductPrice += addProductCountRequest.ProductCount * product.ProductPrice

		err = db.UpdateRecord(&cartProduct, addProductCountRequest.ProductId, "product_id").Error
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
			return
		}

		err = db.FindById(&cart, userId, "user_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
			return
		}

		cart.TotalPrice += cartProduct.ProductPrice
		cart.TotalPrice -= product.ProductPrice
		err = db.UpdateRecord(&cart, cart.CartId, "cart_id").Error
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
			return
		}

		response.ShowResponse(
			"Success",
			utils.HTTP_OK,
			"Product added successfully",
			cartProduct,
			ctx,
		)
		response.ShowResponse(
			"Success",
			utils.HTTP_OK,
			"Cart Updated successfully",
			cart,
			ctx,
		)
	}
}

//Remove all cart service
func RemoveFromCartService(ctx *gin.Context, removeFromCartRequest context.RemoveFromCart) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error in Token header , no userId found")
		return
	}
	if !db.RecordExist("cart_products", "cart_id", removeFromCartRequest.CartId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Cart Id not found")
		return
	}
	if !db.RecordExist("cart_products", "product_id", removeFromCartRequest.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product Id not found")
		return
	}
	var cartProduct model.CartProducts
	var cart model.Cart
	err = db.FindById(&cartProduct, removeFromCartRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error retrieving product from cart_products")
		return
	}

	err = db.FindById(&cart, userId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error retrieving cart details from cart")
		return
	}

	cart.TotalPrice = cart.TotalPrice - cartProduct.ProductPrice
	cart.CartCount = cart.CartCount - 1

	query := "update carts set total_price = 0 , cart_count = 0 where cart_id=?"
	err = db.QueryExecutor(query, &cart, removeFromCartRequest.CartId)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error updating cart")
		return
	}

	err = db.Delete(&cartProduct, removeFromCartRequest.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error Deleting product from cart")
		return
	}
	response.ShowResponse("Success", utils.HTTP_OK, "Product deleted", cartProduct, ctx)
	response.ShowResponse("Success", utils.HTTP_OK, "Cart Details", cart, ctx)
}

// Remove product from cart
func RemoveProductService(ctx *gin.Context, removeProductFromCart context.RemoveProduct) {
	if !db.RecordExist("cart_products", "cart_id", removeProductFromCart.CartId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Cart Id not found")
		return
	}
	if !db.RecordExist("cart_products", "product_id", removeProductFromCart.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product Id not found")
		return
	}

	var product model.Products
	var cartProduct model.CartProducts
	var cart model.Cart
	err := db.FindById(&product, removeProductFromCart.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product Id not found")
		return
	}

	err = db.FindById(&cartProduct, removeProductFromCart.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Cart Id not found")
		return
	}
	if cartProduct.ProductCount < removeProductFromCart.ProductCount {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Request limit exceeded")
		return
	}

	cartProduct.ProductCount -= removeProductFromCart.ProductCount
	cartProduct.ProductPrice -= product.ProductPrice * removeProductFromCart.ProductCount

	if cartProduct.ProductCount <= 0 {

		cartProduct.ProductPrice = 0.0
	}

	err = db.UpdateRecord(&cartProduct, removeProductFromCart.CartId, "cart_id").Error
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error updating cart")
	}

	err = db.FindById(&cart, removeProductFromCart.CartId, "cart_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error finding cart details")
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
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Error updating cart")
	}

	if cartProduct.ProductCount == 0 {
		err = db.Delete(&cartProduct, removeProductFromCart.ProductId, "product_id")
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, err.Error())
			return
		}
	}

	response.ShowResponse("Success", utils.HTTP_OK, "Successfully decremented the product count", cartProduct, ctx)
	response.ShowResponse("Success", utils.HTTP_OK, "Cart Details after decrement", cart, ctx)
}

//Show the cart product details
func GetCartDetailsService(ctx *gin.Context) {
	userId, err := IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Error in Token header , no userId found")
		return
	}

	var cartProductDetails []model.CartProducts
	err = db.FindById(&cartProductDetails, userId, "user_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "No cart products matching this user id")
		return
	}
	var cartResponse []response.CartProductResponse

	for _, product := range cartProductDetails {

		var response response.CartProductResponse
		response.CartId = product.CartId
		response.ProductId = product.ProductId
		response.ProductCount = product.ProductCount
		response.ProductPrice = product.ProductPrice
		response.ProductAddedAt = product.CreatedAt
		cartResponse = append(cartResponse, response)

	}

	response.ShowResponse(
		"Success",
		utils.HTTP_OK,
		"User cart products details are shown below",
		cartResponse,
		ctx,
	)
}
