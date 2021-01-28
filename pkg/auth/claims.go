package auth

import (
	jvt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jvt.StandardClaims
	UserId string `json:"userId"`
	UserIP string `json:"userIP"`
	StartedAt string `json:"startedAt"`
}