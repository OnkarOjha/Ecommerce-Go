package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserId    string `json:"userId"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
