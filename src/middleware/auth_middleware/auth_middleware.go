package middleware

import (
	"strings"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
)

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

		jwtService := service.NewJWTService(nil)

		if jwtService.IsBlacklisted(token) {
			c.Error(&model.AuthenticationError{
				Title:  "Invalid authorization",
				Detail: "The provided token has expired after logging out",
			})

			c.Abort()
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
