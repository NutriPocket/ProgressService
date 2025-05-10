package controller

import (
	"fmt"
	"net/http"

	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/gin-gonic/gin"
)

type FixedDataController struct {
	s service.IUserDataService
}

func NewFixedDataController(s service.IUserDataService) (*FixedDataController, error) {
	var err error

	if s == nil {
		s, err = service.NewUserDataService(nil, nil)
		if err != nil {
			return nil, err
		}
	}

	return &FixedDataController{
		s: s,
	}, nil
}

func (c *FixedDataController) PutFixedData(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	var data *model.BaseFixedUserData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return &model.ValidationError{
			Title:  "Invalid fixed user data",
			Detail: fmt.Sprintf("The fixed user data is invalid, %v", err),
		}
	}

	log.Debugf("Received fixed user data: %v", data)

	data.UserID = authUser.ID
	ret, err, created := c.s.PutFixedData(data)
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

func (c *FixedDataController) GetFixedDataByUser(ctx *gin.Context) error {
	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	jsonRet := make(map[string]any)

	baseData := ctx.Query("base") == "true"

	if baseData {
		var data model.BaseFixedUserData
		data, err = c.s.GetBaseFixedUserDataByUser(authUser.ID)
		jsonRet["data"] = data

	} else {
		var data model.FixedUserData
		data, err = c.s.GetFixedDataByUser(authUser.ID)
		jsonRet["data"] = data
	}

	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, jsonRet)

	return nil
}
