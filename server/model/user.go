package model

import (
	"time"

	"gorm.io/gorm"
)

//DB user representation
type User struct {
	UserId    string         `json:"userId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserName  string         `json:"userName"`
	Gender    string         `json:"gender"`
	Contact   string         `json:"contact"`
	Is_Active bool           `json:"isActive"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
