package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTokenFromAuthHeader(context *gin.Context) (string, error) {

	if context.Request.Header["Authorization"] == nil {
		return "", errors.New("missing authorization header")
	}
	token := strings.Split(context.Request.Header["Authorization"][0], " ")[1]
	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}
