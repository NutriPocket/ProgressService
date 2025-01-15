package routes

import (
	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	{
		users_routes := router.Group("/users")
		users_routes.GET("/", getUsers)
		users_routes.GET("/:username", getUser)
	}
}

func getUsers(c *gin.Context) {}

func getUser(c *gin.Context) {}
