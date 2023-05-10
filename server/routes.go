package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {
	// user handlers
	server.engine.POST("/register", handler.UserRegisterHandler)
	server.engine.POST("/login", handler.UserLoginHandler)
	server.engine.POST("/verify-otp", handler.UserVerifyOtpHandler)
	server.engine.GET("/get-user", gateway.UserAuthorization, handler.GetUserByIdHandler)
	server.engine.POST("/edit-user", gateway.UserAuthorization, handler.EditUserProfile)
	server.engine.DELETE("/logout", gateway.UserAuthorization, handler.UserLogoutHandler)

	// cart handlers
	server.engine.POST("/add-to-cart", gateway.UserAuthorization, handler.AddToCartHandler)
	server.engine.PUT("/add-product", gateway.UserAuthorization, handler.AddProductHandler)
	server.engine.DELETE("/remove-from-cart", gateway.UserAuthorization, handler.RemoveFromCartHandler)
	server.engine.DELETE("/remove-product", gateway.UserAuthorization, handler.RemoveProductHandler)
	server.engine.GET("/get-cart-details", gateway.UserAuthorization, handler.GetCartDetailsHandler)

	//payment handler
	server.engine.POST("/payment", handler.MakePaymentHandler)
	server.engine.GET("/order-details", handler.GetOrderDetails)
	server.engine.PUT("/cancel-order", handler.CancelOrderHandler)
	server.engine.POST("/cart-payment", handler.MakeCartPaymentHandler)

	//filter & search handler
	server.engine.GET("/filter/category", handler.FilterByCategoryHandler)
	server.engine.GET("/filter/price", handler.FilterByPriceHandler)
	server.engine.GET("/filter/brand", handler.FilterByBrandHandler)
	server.engine.GET("/search-bar", handler.SearchBarHandler)
	server.engine.GET("/searchbar/history", handler.SearchBarHistoryHandler)

	//user address handler
	server.engine.POST("/user/address", handler.UserAddressHandler)
	server.engine.GET("user/address-get", handler.UserAddressRetrieveHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
