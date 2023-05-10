package gateway

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func UserAuthorization(c *gin.Context) {

	tokenString, err := utils.GetTokenFromAuthHeader(c)
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Error fetching token",
		)
		c.Abort()
		return

	}
	claims, err := token.DecodeToken(c, tokenString)
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Error :"+err.Error(),
		)
		c.Abort()
		return
	}

	if claims.Role != "customer" {
		response.ErrorResponse(c, utils.HTTP_UNAUTHORIZED, "Token doesnot belong to a customer")
		return
	}

	var userSession model.Session
	err = db.FindById(&userSession, claims.UserId, "user_id")
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "User Id retrieved through token does not exist",
		)
		c.Abort()
		return
	}

	err = db.FindById(&userSession, tokenString, "token")
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Token does not match any session",
		)
		c.Abort()
		return
	}

	c.Next()

}

func VendorAuthorization(c *gin.Context) {

	tokenString, err := utils.GetTokenFromAuthHeader(c)
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Error fetching token",
		)
		c.Abort()
		return

	}
	claims, err := token.DecodeToken(c, tokenString)
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Error :"+err.Error(),
		)
		c.Abort()
		return
	}

	if claims.Role != "vendor" {
		response.ErrorResponse(c, utils.HTTP_UNAUTHORIZED, "Token doesnot belong to a vendor")
		return
	}

	var userSession model.Session
	err = db.FindById(&userSession, claims.UserId, "user_id")
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "User Id retrieved through token does not exist",
		)
		c.Abort()
		return
	}

	err = db.FindById(&userSession, tokenString, "token")
	if err != nil {
		response.ErrorResponse(
			c, utils.HTTP_UNAUTHORIZED, "Token does not match any session",
		)
		c.Abort()
		return
	}

	c.Next()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(utils.HTTP_NO_CONTENT)
			return
		}

		c.Next()
	}
}
