package model

import "gorm.io/gorm"

type UserAddresses struct {
	gorm.Model
	UserId string `json:"userId"`
	Address string `json:"address"`
	AddressType string `json:"addressType"`
}
