// Package model contains the structs types that will be used in the application.
package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWTPayload is a struct that contains the User data and the JWT claims
type JWTPayload struct {
	// Payload is the User data
	Payload User `json:"payload"`
	jwt.RegisteredClaims
}
