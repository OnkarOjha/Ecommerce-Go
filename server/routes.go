package server

import (
	_ "main/docs"
	"main/server/handler"
	"main/server/provider"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	// user handlers
	server.engine.POST("/register", handler.UserRegisterHandler)
	server.engine.POST("/login", handler.UserLoginHandler)
	server.engine.POST("/verifyOtp", handler.UserVerifyOtpHandler)
	server.engine.POST("/getUser", handler.GetUserByIdHandler)
	server.engine.POST("/editUser", provider.UserAuthorization, handler.EditUserProfile)
	server.engine.DELETE("/logout", provider.UserAuthorization, handler.UserLogoutHandler)

	// cart handlers
	server.engine.POST("/addToCart", provider.UserAuthorization, handler.AddToCartHandler)
	server.engine.PUT("/addProduct", provider.UserAuthorization, handler.AddProductHandler)
	server.engine.DELETE("/removeFromCart", provider.UserAuthorization, handler.RemoveFromCartHandler)
	server.engine.DELETE("/removeProduct", provider.UserAuthorization, handler.RemoveProductHandler)
	server.engine.GET("/getCartDetails", provider.UserAuthorization, handler.GetCartDetailsHandler)

	//payment handler
	server.engine.POST("/payment", handler.MakePaymentHandler)
	server.engine.GET("/orderDetails", handler.GetOrderDetails)

	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
