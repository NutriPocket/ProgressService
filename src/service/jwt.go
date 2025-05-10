// Package service contains the services that will be used in the application.
package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/NutriPocket/UserService/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService is a struct that will be used to sign, verify, decode and blacklist JWT tokens.
type JWTService struct {
	// key is the secret key used to sign the JWT tokens.
	key []byte
}

var jwtKey = os.Getenv("JWT_SECRET_KEY")

// NewJWTService creates a new JWTService with the provided IJWTRepository.
// jwtRepository is the repository that will be used to interact with the jwt_blacklist table.
// It returns a new JWTService.
func NewJWTService() (*JWTService, error) {
	var key = []byte("secret")

	if jwtKey != "" {
		key = []byte(jwtKey)
	}

	return &JWTService{key: []byte(key)}, nil
}

// Sign signs a JWT token with the provided payload.
// payload is the user data to sign.
// It returns the signed token and an error if the operation fails.
func (service *JWTService) Sign(payload model.User) (string, error) {
	nowUtc := time.Now().UTC()

	claim := model.JWTPayload{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(nowUtc.Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(nowUtc),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(service.key)

	return tokenString, err
}

// isJWT checks if a token has the JWT format.
// tokenString is the token to check.
// It returns true if the token has the JWT format, false otherwise.
func (service *JWTService) isJWT(tokenString string) bool {
	jwtRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)$`)

	return jwtRegex.MatchString(tokenString)
}

// Verify verifies a JWT token.
// tokenString is the token to verify.
// It returns true if the token is valid, false otherwise.
func (service *JWTService) Verify(tokenString string) (bool, error) {
	if !service.isJWT(tokenString) {
		return false, &model.ValidationError{Title: "Invalid JWT", Detail: "The provided token doesn't have JWT format"}
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return service.key, nil
	})

	return token.Valid, err
}

// Decode decodes a JWT token.
// tokenString is the token to decode.
// It returns the decoded payload and an error if the operation fails.
func (service *JWTService) Decode(tokenString string) (model.JWTPayload, error) {
	if !service.isJWT(tokenString) {
		return model.JWTPayload{}, &model.ValidationError{
			Title:  "Invalid JWT",
			Detail: "The provided token doesn't have JWT format",
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return service.key, nil
	})

	if claims, ok := token.Claims.(*model.JWTPayload); ok && token.Valid {
		return *claims, nil
	} else {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return model.JWTPayload{}, &model.AuthenticationError{
				Title:  "Expired token",
				Detail: "Your token has expired, please try logging in again",
			}
		}

		return model.JWTPayload{}, err
	}
}
