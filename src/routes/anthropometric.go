// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"github.com/NutriPocket/ProgressService/controller"
	"github.com/gin-gonic/gin"
)

func putAnthropometricData(c *gin.Context) {
	controller, err := controller.NewAnthropometricController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.PutAnthropometricData(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func getAnthropometricData(c *gin.Context) {
	controller, err := controller.NewAnthropometricController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetAnthropometricDataByUser(c)
	if err != nil {
		c.Error(err)
		return
	}
}
