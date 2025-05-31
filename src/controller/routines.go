package controller

import (
	"fmt"
	"net/http"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/service"
	"github.com/gin-gonic/gin"
)

type RoutineController struct {
	s service.IRoutineService
}

func NewRoutineController(s service.IRoutineService) (*RoutineController, error) {
	var err error

	if s == nil {
		s, err = service.NewRoutineService(nil)
		if err != nil {
			return nil, err
		}
	}

	return &RoutineController{
		s: s,
	}, nil
}

func (c *RoutineController) PostRoutine(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.RoutineDTO
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid routine data",
			Detail: fmt.Sprintf("The routine data is invalid, %v", err),
		}
	}

	data.UserID = authUser.ID
	ret, err := c.s.CreateRoutine(data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = ret

	ctx.JSON(http.StatusCreated, jsonRet)
	return nil
}

func (c *RoutineController) GetRoutinesByUser(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data []model.RoutineData
	err = c.s.GetRoutinesByUser(authUser.ID, &data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = data

	ctx.JSON(http.StatusOK, jsonRet)
	return nil
}

func (c *RoutineController) GetFreeSchedules(ctx *gin.Context) error {
	users := ctx.QueryArray("users")

	data, err := c.s.GetFreeSchedules(users)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = data

	ctx.JSON(http.StatusOK, jsonRet)
	return nil
}

func (c *RoutineController) DeleteRoutineBySchedule(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.Schedule
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid routine data",
			Detail: fmt.Sprintf("The routine data is invalid, %v", err),
		}
	}

	ret, err := c.s.DeleteRutineBySchedule(authUser.ID, data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = ret

	ctx.JSON(http.StatusOK, jsonRet)
	return nil
}
