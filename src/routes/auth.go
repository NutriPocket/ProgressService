// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"net/http"

	controller "github.com/NutriPocket/UserService/controller/users"
	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	{
		auth_routes := router.Group("/auth")
		auth_routes.POST("/register", register)
		auth_routes.POST("/login", login)
		auth_routes.POST("/logout", logout)
	}
}

func register(c *gin.Context) {
	var userData model.BaseUser

	if err := c.BindJSON(&userData); err != nil {
		c.Error(&model.ValidationError{
			Title:  "Wrong body format",
			Detail: "Expected a json body with the user account data in it",
		})
		return
	}

	controller := controller.UserController{}

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

	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	service, err := service.NewUserService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	createdUser, err := service.CreateUser(&userData)

	if err != nil {
		c.Error(err)
		return
	}

	signed, err := jwtService.Sign(createdUser)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdUser, "token": signed})
}

func login(c *gin.Context) {
	body := model.LoginUser{}

	if err := c.BindJSON(&body); err != nil {
		c.Error(&model.ValidationError{
			Title:  "Wrong body format",
			Detail: "Expected a json body with the user credentials in it",
		})
		return
	}

	controller := controller.UserController{}

	controller.ValidateUsernameOrEmail(body.EmailOrUsername)
	controller.ValidateString(body.Password, "password")

	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	service, err := service.NewUserService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := service.Login(&body)

	if err != nil {
		c.Error(err)
		return
	}

	signed, err := jwtService.Sign(user)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "token": signed})
}

func logout(c *gin.Context) {
	body := struct{ Token string }{}

	if err := c.BindJSON(&body); err != nil {
		c.Error(&model.ValidationError{
			Title:  "Wrong body format",
			Detail: "Expected a json body with the key 'token' in it",
		})
		return
	}

	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	if err := jwtService.Blacklist(body.Token); err != nil {
		c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
