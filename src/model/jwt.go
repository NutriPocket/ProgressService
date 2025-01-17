package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	Payload User `json:"payload"`
	jwt.RegisteredClaims
}
