// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"github.com/NutriPocket/ProgressService/controller"
	"github.com/gin-gonic/gin"
)

func putFixedData(c *gin.Context) {
	controller, err := controller.NewFixedDataController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.PutFixedData(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func getFixedData(c *gin.Context) {
	controller, err := controller.NewFixedDataController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetFixedDataByUser(c)

	if err != nil {
		c.Error(err)
		return
	}
}
