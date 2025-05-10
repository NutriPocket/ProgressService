package controller

import (
	"time"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/gin-gonic/gin"
)

var authError = &model.AuthenticationError{
	Title:  "Unauthorized user",
	Detail: "The user isn't authorized to access this endpoint",
}

func GetAuthUser(c *gin.Context) (*model.User, error) {
	var user *model.User
	var ok bool

	authUser, exists := c.Get("authUser")
	if !exists {
		return nil, authError
	}

	if user, ok = authUser.(*model.User); !ok {
		return nil, authError
	}

	if userId := c.Param("userId"); userId != user.ID {
		return nil, authError
	}

	return user, nil
}

func ValidateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return &model.ValidationError{
			Title:  "Invalid date",
			Detail: "The format of the date is invalid, expected format: YYYY-MM-DD",
		}
	}

	return nil
}

func ValidateDeadline(date string) error {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return &model.ValidationError{
			Title:  "Invalid date",
			Detail: "The format of the date is invalid, expected format: YYYY-MM-DD",
		}
	}

	if parsed.Before(time.Now()) {
		return &model.ValidationError{
			Title:  "Invalid deadline",
			Detail: "The deadline must be in the future",
		}
	}

	return nil
}
