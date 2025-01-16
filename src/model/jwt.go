package model

import "github.com/golang-jwt/jwt/v5"

type JWTPayload struct {
	Payload interface{} `json:"payload"`
	jwt.RegisteredClaims
}
