package model

import "time"

// DB model to store session information
type Session struct {
	SessionId string    `json:"sessionId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId    string    `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
