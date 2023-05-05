package model

type Session struct {
	SessionId string `json:"sessionId" gorm:"default:uuid_generate_v4();unique;primaryKey"`
	UserId    string `json:"userId"`
	Token     string `json:"token"`
}
