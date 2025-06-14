package routes

import (
	"github.com/NutriPocket/ProgressService/controller"
	"github.com/gin-gonic/gin"
)

func postExercise(c *gin.Context) {
	controller, err := controller.NewExerciseController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.CreateExercise(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func getExercisesByUserIdAndDate(c *gin.Context) {
	controller, err := controller.NewExerciseController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.GetExercisesByUser(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func putExercise(c *gin.Context) {
	controller, err := controller.NewExerciseController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.UpdateExercise(c)
	if err != nil {
		c.Error(err)
		return
	}
}

func deleteExercise(c *gin.Context) {
	controller, err := controller.NewExerciseController(nil)
	if err != nil {
		c.Error(err)
		return
	}

	err = controller.DeleteExercise(c)
	if err != nil {
		c.Error(err)
		return
	}
}
