package model

import "gorm.io/gorm"

type DbConstant struct {
	gorm.Model
	ConstantName      string `json:"constantName"`
	ConstantShortHand string `json:"constantShortHand"`
}
