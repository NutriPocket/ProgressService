package utils

import (
	"github.com/MaxiOtero6/go-auth-rest/routes"
	"github.com/gin-gonic/gin"

	middlewareAuth "github.com/MaxiOtero6/go-auth-rest/middleware/auth_middleware"
	middlewareErr "github.com/MaxiOtero6/go-auth-rest/middleware/error_handler"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middlewareErr.ErrorHandler())
	router.Use(middlewareAuth.AuthMiddleware())
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)

	return router
}
