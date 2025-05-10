package controller

import (
	"fmt"
	"net/http"

	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type AnthropometricController struct {
	s service.IUserDataService
}

func NewAnthropometricController(s service.IUserDataService) (*AnthropometricController, error) {
	var err error

	if s == nil {
		s, err = service.NewUserDataService(nil, nil)
		if err != nil {
			return nil, err
		}
	}

	return &AnthropometricController{
		s: s,
	}, nil
}

func (c *AnthropometricController) PutAnthropometricData(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.AnthropometricData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid anthropometric user data",
			Detail: fmt.Sprintf("The anthropometric user data is invalid, %v", err),
		}
	}

	log.Debugf("Received anthropometric data: %v", data)

	data.UserID = authUser.ID
	ret, err, created := c.s.PutAnthropometricData(data)
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

func (c *AnthropometricController) GetAnthropometricDataByUser(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)

	date := ctx.Query("date")

	if date != "" {
		if err := ValidateDate(date); err != nil {
			return err
		}

		var data model.AnthropometricData
		data, err = c.s.GetAnthropometricDataByUserAndDay(authUser.ID, date)
		jsonRet["data"] = data
	} else {
		var data []model.AnthropometricData
		data, err = c.s.GetAllAnthropometricDataByUser(authUser.ID)
		jsonRet["data"] = data
	}

	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, jsonRet)

	return nil
}
