package routes

import (
	"net/http"

	"github.com/MaxiOtero6/go-auth-rest/controller"
	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	{
		auth_routes := router.Group("/auth")
		auth_routes.POST("/register", register)
		auth_routes.POST("/login", login)
		auth_routes.POST("/logout", logout)
		auth_routes.POST("/refresh", refresh)
	}
}

func register(c *gin.Context) {
	var userData model.BaseUser

	if err := c.BindJSON(&userData); err != nil {
		c.Error(err)
		return
	}

	if err := controller.ValidateString(userData.Username, "username"); err != nil {
		c.Error(err)
		return
	}

	if err := controller.ValidateString(userData.Password, "password"); err != nil {
		c.Error(err)
		return
	}

	if err := controller.ValidateEmail(userData.Email); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, service.CreateUser(&userData))
}

func login(c *gin.Context) {}

func logout(c *gin.Context) {}

func refresh(c *gin.Context) {}
