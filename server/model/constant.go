package model

import "gorm.io/gorm"

// DB model to represent DB constants
type DbConstant struct {
	gorm.Model
	ConstantName      string `json:"constantName"`
	ConstantShortHand string `json:"constantShortHand"`
}
