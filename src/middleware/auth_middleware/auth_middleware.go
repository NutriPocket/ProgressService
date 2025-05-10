// Package middleware provides custom middlewares for the API
package middleware

import (
	"strings"

	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/gin-gonic/gin"
)

// getRootPath returns the root path of a URL
// @param urlPath string - The URL path
func getRootPath(urlPath string) string {
	if urlPath == "" {
		return ""
	}

	urlPath = urlPath[1:]

	if idx := strings.Index(urlPath, "/"); idx != -1 {
		return urlPath[:idx]
	} else {
		return urlPath
	}
}

// getToken returns the token from the Authorization header
// authHeader is the Authorization header
// It returns the token parsed from the Authorization header
// It returns an error if the Authorization header is empty or has an invalid format
func getToken(authHeader string) (token string, err error) {
	if authHeader == "" {
		return "", &model.AuthenticationError{
			Title:  "Unauthorized user",
			Detail: "The user isn't authorized because no Authorization header is provided",
		}
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", &model.AuthenticationError{
			Title:  "Invalid authorization",
			Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
		}
	}

	token = parts[1]

	return token, nil
}

// AuthMiddleware is a middleware that checks if the user is authorized to access the endpoint
// Only the endpoints that start with /auth are allowed to be accessed without authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path

		if getRootPath(urlPath) == "auth" {
			c.Next()

			return
		}

		authHeader := c.GetHeader("Authorization")

		token, err := getToken(authHeader)

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		jwtService, err := service.NewJWTService()
		if err != nil {
			c.Error(err)
			return
		}

		decoded, err := jwtService.Decode(token)

		if err != nil {
			c.Error(err)

			c.Abort()
			return
		}

		c.Set("authUser", decoded.Payload)

		c.Next()
	}
}
