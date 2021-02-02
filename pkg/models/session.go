package models

type Session struct {
	UserId string `json:"userId" bson:"_id"`
	UserIP string `json:"userIP" bson:"userIP"`
	StartedAt string `json:"startedAt,omitempty" bson:"startedAt"`
	RefreshToken string `json:"refreshToken,omitempty" bson:"refreshToken"`
}