package routes

import (
	"github.com/NutriPocket/ProgressService/controller"
	"github.com/gin-gonic/gin"
)

func postRoutine(c *gin.Context) {
	controller, err := controller.NewRoutineController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.PostRoutine(c)

	if err != nil {
		c.Error(err)
		return
	}
}

func getRoutines(c *gin.Context) {
	controller, err := controller.NewRoutineController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetRoutinesByUser(c)

	if err != nil {
		c.Error(err)
		return
	}
}

func getFreeSchedules(c *gin.Context) {
	controller, err := controller.NewRoutineController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetFreeSchedules(c)

	if err != nil {
		c.Error(err)
		return
	}
}

func deleteRoutine(c *gin.Context) {
	controller, err := controller.NewRoutineController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.DeleteRoutineBySchedule(c)

	if err != nil {
		c.Error(err)
		return
	}
}
