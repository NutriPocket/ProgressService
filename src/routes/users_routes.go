package routes

import "github.com/gin-gonic/gin"

func UsersRoutes(router *gin.Engine) {
	{
		routes := router.Group("/users")
		/*
			Anthropometric Data routes
		*/
		routes.PUT("/:userId/anthropometrics/", putAnthropometricData)
		routes.GET("/:userId/anthropometrics/", getAnthropometricData)
		/*
			Fixed User Data routes
		*/
		routes.PUT("/:userId/fixedData/", putFixedData)
		routes.GET("/:userId/fixedData/", getFixedData)
		/*
			Objectives routes
		*/
		routes.PUT("/:userId/objectives/", putObjectiveData)
		routes.GET("/:userId/objectives/", getObjectiveData)
		/*
			Routines routes
		*/
		routes.POST("/:userId/routines/", postRoutine)
		routes.GET("/:userId/routines/", getRoutines)
		routes.DELETE("/:userId/routines/", deleteRoutine)
		routes.GET("/freeSchedules/", getFreeSchedules)
		/*
			Exercise routes
		*/
		routes.POST("/:userId/exercises/", postExercise)
		routes.GET("/:userId/exercises/", getExercisesByUserIdAndDate)
		routes.PUT("/:userId/exercises/:id", putExercise)
		routes.DELETE("/:userId/exercises/:id", deleteExercise)
	}
}
