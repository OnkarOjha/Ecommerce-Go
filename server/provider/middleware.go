package provider

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func UserAuthorization(c *gin.Context) {

	// 1. user id exists in session table
	// 2. token associated with userId is same or not
	// 3. token validity
	tokenString, err := utils.GetTokenFromAuthHeader(c)
	if err != nil {
		response.ErrorResponse(
			c, 401, "Error fetching token",
		)
		c.Abort()
		return
	}
	claims, err := DecodeToken(c, tokenString)
	if err != nil {
		response.ErrorResponse(
			c, 401, "Error :"+err.Error(),
		)
		c.Abort()
		return
	}

	var userSession model.Session
	err = db.FindById(&userSession, claims.UserId, "user_id")
	if err != nil {
		response.ErrorResponse(
			c, 401, "User Id retrieved through token does not exist",
		)
		c.Abort()
		return
	}

	err = db.FindById(&userSession, tokenString, "token")
	if err != nil {
		response.ErrorResponse(
			c, 400, "Token does not match any session",
		)
		c.Abort()
	}

	c.Next()

}

func AdminAuthorization(c *gin.Context) {

	// fmt.Println("inside middleware")
	// token, err := c.Request.Cookie("cookie")
	// if err != nil {

	// 	response.ErrorResponse(c, 400, err.Error())
	// 	c.Abort()
	// 	return

	// }

	// claims, err := DecodeToken(token.Value)
	// if err != nil {
	// 	response.ErrorResponse(c, 401, err.Error())
	// 	c.Abort()
	// 	return
	// }
	// err = claims.Valid()
	// if err != nil {
	// 	response.ErrorResponse(c, 401, err.Error())
	// 	c.Abort()
	// 	return
	// }
	// if claims.Role == "admin" {
	// 	c.Next()
	// } else {
	// 	response.ErrorResponse(c, 403, "Access Denied")
	// 	c.Abort()
	// 	return

	// }

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
