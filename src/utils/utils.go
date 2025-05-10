// Package utils provide utility functions that are used in the application.
package utils

import (
	"github.com/NutriPocket/ProgressService/routes"
	"github.com/gin-gonic/gin"

	middlewareAuth "github.com/NutriPocket/ProgressService/middleware/auth_middleware"
	middlewareErr "github.com/NutriPocket/ProgressService/middleware/error_handler"
)

// SetupRouter sets up the routes for the application.
// It returns a router with the middlewares and routes set up.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middlewareErr.ErrorHandler())
	router.Use(middlewareAuth.AuthMiddleware())
	routes.AnthropometricRoutes(router)
	routes.FixedUserDataRoutes(router)

	return router
}
