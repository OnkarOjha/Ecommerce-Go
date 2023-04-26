package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId string `json:"userId"`
	UserName string `json:"userName"`
	Gender string `json:"gender"`
	Contact string `json:"contact"`
	Is_Active bool `json:"isActive"`
}
