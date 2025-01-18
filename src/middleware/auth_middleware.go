package middleware

import (
	"strings"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path[1:]

		var route string
		if idx := strings.Index(urlPath, "/"); idx != -1 {
			route = urlPath[:idx]
		} else {
			route = urlPath
		}

		if route == "auth" {
			c.Next()

			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(&model.AuthenticationError{
				Title:  "Unauthorized user",
				Detail: "The user isn't authorized because no Authorization header is provided",
			})

			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(&model.AuthenticationError{
				Title:  "Invalid authorization",
				Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
			})

			c.Abort()
			return
		}

		token := parts[1]

		jwtService := service.NewJWTService()

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
