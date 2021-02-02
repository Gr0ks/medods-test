package models

type Session struct {
	UserId string `json:"userId"`
	UserIP string `json:"userIP"`
	StartedAt string `json:"startedAt,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}