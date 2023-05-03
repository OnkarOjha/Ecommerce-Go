package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTokenFromAuthHeader(context *gin.Context) (string, error) {
	token := strings.Split(context.Request.Header["Authorization"][0], " ")[1]
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}
