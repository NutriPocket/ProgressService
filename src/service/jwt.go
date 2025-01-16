package service

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	key []byte
}

var jwtKey = os.Getenv("JWT_SECRET_KEY")

func NewJWTService() JWTService {
	var key = []byte("secret")

	if jwtKey != "" {
		key = []byte(jwtKey)
	}

	return JWTService{key: []byte(key)}
}

func (service *JWTService) Sign(payload interface{}) (string, error) {
	claim := model.JWTPayload{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate((time.Now().Add(time.Minute * 5))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(service.key)

	return tokenString, err
}

func (service *JWTService) isJWT(tokenString string) bool {
	jwtRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)$`)

	return jwtRegex.MatchString(tokenString)
}

func (service *JWTService) Verify(tokenString string) (bool, error) {
	if !service.isJWT(tokenString) {
		return false, fmt.Errorf("the provided token is not a jwt string")
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return service.key, nil
	})

	return token.Valid, err
}

func (service *JWTService) Decode(tokenString string) (model.JWTPayload, error) {
	if !service.isJWT(tokenString) {
		return model.JWTPayload{}, fmt.Errorf("the provided token is not a jwt string")
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
		return model.JWTPayload{}, err
	}
}
