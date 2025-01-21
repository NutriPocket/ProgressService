// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"net/http"

	controller "github.com/MaxiOtero6/go-auth-rest/controller/users"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	{
		users_routes := router.Group("/users")
		users_routes.GET("/", getUsers)
		users_routes.GET("/:username", getUser)
	}
}

func getUsers(c *gin.Context) {
	service := service.NewUserService(nil)

	users, err := service.GetAllUsers()

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	username := c.Param("username")

	controller := controller.UserController{}

	controller.ValidateString(username, "username")

	service := service.NewUserService(nil)

	user, err := service.GetUser(username)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}
