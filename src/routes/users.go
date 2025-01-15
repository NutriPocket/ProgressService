package routes

import (
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersRoutes(router *gin.Engine) {
	{
		users_routes := router.Group("/users")
		users_routes.GET("/", getUsers)
		users_routes.GET("/:username", getUser)
	}
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllUsers())
}

func getUser(c *gin.Context) {}
