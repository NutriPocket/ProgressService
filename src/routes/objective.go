// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"github.com/NutriPocket/ProgressService/controller"
	"github.com/gin-gonic/gin"
)

func putObjectiveData(c *gin.Context) {
	controller, err := controller.NewObjectiveController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.PutObjective(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func getObjectiveData(c *gin.Context) {
	controller, err := controller.NewObjectiveController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetObjectiveByUser(c)
	if err != nil {
		c.Error(err)
		return
	}
}
