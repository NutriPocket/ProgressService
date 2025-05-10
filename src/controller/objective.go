package controller

import (
	"fmt"
	"net/http"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/service"
	"github.com/gin-gonic/gin"
)

type ObjectiveController struct {
	s service.IObjectiveService
}

func NewObjectiveController(s service.IObjectiveService) (*ObjectiveController, error) {
	var err error

	if s == nil {
		s, err = service.NewObjectiveService(nil)
		if err != nil {
			return nil, err
		}
	}

	return &ObjectiveController{
		s: s,
	}, nil
}

func (c *ObjectiveController) PutObjective(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.ObjectiveData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid user objective data",
			Detail: fmt.Sprintf("The user objective data is invalid, %v", err),
		}
	}

	if err := ValidateDeadline(data.Deadline); err != nil {
		return err
	}

	log.Debugf("Received user objective data: %v", data)

	data.UserID = authUser.ID
	ret, err, created := c.s.PutObjective(data)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)
	jsonRet["data"] = ret

	if created {
		ctx.JSON(http.StatusCreated, jsonRet)
	} else {
		ctx.JSON(http.StatusOK, jsonRet)
	}

	return nil
}

func (c *ObjectiveController) GetObjectiveByUser(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)

	var data model.ObjectiveData
	data, err = c.s.GetObjectiveByUser(authUser.ID)
	jsonRet["data"] = data

	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, jsonRet)

	return nil
}
