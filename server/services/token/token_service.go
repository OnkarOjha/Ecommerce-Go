package token

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId string `json:"userId"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

//Generate JWT Token
func GenerateToken(claims Claims, context *gin.Context) string {
	//create user claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))

	if err != nil {
		response.ErrorResponse(context, utils.HTTP_UNAUTHORIZED, "Error signing token")
	}
	return tokenString
}

//Decode Token function
func DecodeToken(context *gin.Context, tokenString string) (Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		if claims.ExpiresAt != nil && (*claims.ExpiresAt).Before(time.Now()) {
			var userToBeLoggedOut model.User
			err := db.FindById(&userToBeLoggedOut, claims.UserId, "user_id")
			if err != nil {
				return *claims, fmt.Errorf("error finding user in db")
			}
			query := "UPDATE users SET is_active = false WHERE user_id = '" + claims.UserId + "'"
			db.QueryExecutor(query, &userToBeLoggedOut)
			fmt.Println("user to be logged out ", userToBeLoggedOut)

			var userSessionToBeDeleted model.Session
			err = db.FindById(&userSessionToBeDeleted, claims.UserId, "user_id")
			if err != nil {
				return *claims, fmt.Errorf("error finding user in db")
			}

			db.DeleteRecord(&userSessionToBeDeleted, claims.UserId, "user_id")

			return *claims, fmt.Errorf("token has expired , please proceed to login")
		}
		return *claims, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
