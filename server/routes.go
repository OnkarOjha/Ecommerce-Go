package server

import (
	_ "main/docs"
	"main/server/gateway"
	"main/server/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {
	// use cors middleware
	server.engine.Use(gateway.CORSMiddleware())

	// coupon expiration
	server.engine.Use(gateway.CouponExpiration)

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
	server.engine.POST("/payment", gateway.UserAuthorization, handler.MakePaymentHandler)
	server.engine.GET("/order-details", gateway.UserAuthorization, handler.GetOrderDetails)
	server.engine.PUT("/cancel-order", gateway.UserAuthorization, handler.CancelOrderHandler)
	server.engine.POST("/cart-payment", gateway.UserAuthorization, handler.MakeCartPaymentHandler)

	//filter & search handler
	server.engine.GET("/filter/category", handler.FilterByCategoryHandler)
	server.engine.GET("/filter/price", handler.FilterByPriceHandler)
	server.engine.GET("/filter/brand", handler.FilterByBrandHandler)
	server.engine.GET("/search-bar", handler.SearchBarHandler)
	server.engine.GET("/search-bar/history", handler.SearchBarHistoryHandler)

	//user address handler
	server.engine.POST("/user/address", handler.UserAddressHandler)
	server.engine.GET("/user/address-get", handler.UserAddressRetrieveHandler)

	//vendor login  & logout
	server.engine.POST("/vendor-register", handler.VendorRegisterHandler)
	server.engine.POST("/vendor-login", handler.VendorLoginHandler)
	server.engine.POST("/vendor-verify-otp", handler.VenderVerifyOtpHandler)
	server.engine.DELETE("/vendor-logout", gateway.VendorAuthorization, handler.VendorLogoutHandler)
	server.engine.POST("/vendor-edit-details", gateway.VendorAuthorization, handler.VendorEditDetailsHandler)
	server.engine.POST("/vendor-file-upload", gateway.VendorAuthorization, handler.VendorFileUpload)
	server.engine.GET("/file-get", gateway.VendorAuthorization, handler.VendorFileGet)

	//vendor product management
	server.engine.POST("/vendor-product-add", gateway.VendorAuthorization, handler.InventoryProductAddHandler)
	server.engine.POST("/vendor-product-update", gateway.VendorAuthorization, handler.InventoryProductUpdateHandler)
	server.engine.POST("/vendor-product-delete", gateway.VendorAuthorization, handler.InventoryProductDeleteHandler)
	server.engine.POST("/vendor-product-status", gateway.VendorAuthorization, handler.VendorOrderStatusUpdateHandler)

	//coupon codes
	server.engine.POST("/coupon-add", handler.AddCouponHandler)
	server.engine.GET("/coupon-get", handler.GetCouponsHandler)

	// common routes
	server.engine.POST("/user-delete", handler.UserDeleteHandler)
	server.engine.GET("/username-exists", handler.UserNameHandler)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
