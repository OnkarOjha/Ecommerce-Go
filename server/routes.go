package server

import (
	_ "main/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConfigureRoutes(server *Server) {
	
	server.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
}
