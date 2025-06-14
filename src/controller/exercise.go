package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/service"
	"github.com/gin-gonic/gin"
)

type ExerciseController struct {
	s service.IExerciseService
}

func NewExerciseController(s service.IExerciseService) (*ExerciseController, error) {
	var err error

	if s == nil {
		s, err = service.NewExerciseService(nil)
		if err != nil {
			return nil, err
		}
	}

	return &ExerciseController{
		s: s,
	}, nil
}

// CreateExercise handles POST requests to create a new exercise
func (c *ExerciseController) CreateExercise(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.ExerciseDTO
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid exercise data",
			Detail: fmt.Sprintf("The exercise data is invalid, %v", err),
		}
	}

	log.Debugf("Received exercise data: %v", data)

	data.UserID = authUser.ID
	ret, err := c.s.CreateExercise(data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = ret

	ctx.JSON(http.StatusCreated, jsonRet)
	return nil
}

// GetExercisesByUser handles GET requests to retrieve all exercises for a user
func (c *ExerciseController) GetExercisesByUser(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	// Check if date parameter is provided
	date := ctx.Query("date")
	exercises, err := c.s.GetExercisesByUserIdAndDate(authUser.ID, date)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = exercises

	ctx.JSON(http.StatusOK, jsonRet)
	return nil
}

// UpdateExercise handles PUT requests to update an existing exercise
func (c *ExerciseController) UpdateExercise(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return &model.ValidationError{
			Title:  "Invalid exercise ID",
			Detail: "Exercise ID must be a positive integer",
		}
	}

	var data *model.ExerciseDTO
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid exercise data",
			Detail: fmt.Sprintf("The exercise data is invalid, %v", err),
		}
	}

	// Ensure the user ID cannot be changed
	data.UserID = authUser.ID

	ret, err := c.s.UpdateExercise(id, authUser.ID, data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = ret

	ctx.JSON(http.StatusOK, jsonRet)
	return nil
}

// DeleteExercise handles DELETE requests to remove an exercise
func (c *ExerciseController) DeleteExercise(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return &model.ValidationError{
			Title:  "Invalid exercise ID",
			Detail: "Exercise ID must be a positive integer",
		}
	}

	err = c.s.DeleteExercise(id, authUser.ID)
	if err != nil {
		return err
	}

	ctx.Status(http.StatusNoContent)
	return nil
}
