// Package utils provide utility functions that are used in the application.
package utils

import (
	"github.com/MaxiOtero6/go-auth-rest/routes"
	"github.com/gin-gonic/gin"

	middlewareAuth "github.com/MaxiOtero6/go-auth-rest/middleware/auth_middleware"
	middlewareErr "github.com/MaxiOtero6/go-auth-rest/middleware/error_handler"
)

// SetupRouter sets up the routes for the application.
// It returns a router with the middlewares and routes set up.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middlewareErr.ErrorHandler())
	router.Use(middlewareAuth.AuthMiddleware())
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)

	return router
}
