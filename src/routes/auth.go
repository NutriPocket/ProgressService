package routes

import (
	"net/http"

	controller "github.com/MaxiOtero6/go-auth-rest/controller/users"
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
	}
}

func register(c *gin.Context) {
	var userData model.BaseUser

	if err := c.BindJSON(&userData); err != nil {
		c.Error(err)
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

	jwtService := service.NewJWTService(nil)
	service := service.UserService{}

	createdUser := service.CreateUser(&userData)

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
		c.Error(err)
		return
	}

	controller := controller.UserController{}

	controller.ValidateUsernameOrEmail(body.EmailOrUsername)
	controller.ValidateString(body.Password, "password")

	jwtService := service.NewJWTService(nil)
	service := service.UserService{}

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
		c.Error(err)
		return
	}

	jwtService := service.NewJWTService(nil)

	if err := jwtService.Blacklist(body.Token); err != nil {
		c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
