package server

import (
	_ "main/docs"
	"main/server/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {

	// user handlers
	server.engine.POST("/register" , handler.UserRegisterHandler)
	server.engine.POST("/login" , handler.UserLoginHandler)
	server.engine.POST("/verifyOtp" , handler.UserVerifyOtpHandler)
	server.engine.POST("/getUser" , handler.GetUserByIdHandler)
	server.engine.POST("/editUser" , handler.EditUserProfile)





	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
}
