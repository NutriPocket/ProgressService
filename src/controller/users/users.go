package controller

import (
	"regexp"
	"strings"

	"github.com/MaxiOtero6/go-auth-rest/model"
)

type UserController struct{}

func (controller *UserController) ValidateString(str string, field string) error {
	if str == "" {
		return &model.ValidationError{Detail: "The " + field + " field is required", Title: "Empty " + field + " field"}
	}

	if len(str) > 100 {
		return &model.ValidationError{Detail: "The " + field + " field must be less than 100 characters", Title: "Invalid " + field + " field"}
	}

	return nil
}

func (controller *UserController) ValidateEmail(email string) error {
	var err error

	if err = controller.ValidateString(email, "email"); err != nil {
		return err
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`)

	if !emailRegex.MatchString(email) {
		return &model.ValidationError{Detail: "The email field must be a valid email address", Title: "Invalid email field"}
	}

	return nil
}

func (controller *UserController) ValidateUsernameOrEmail(str string) error {
	if strings.Contains(str, "@") {
		return controller.ValidateEmail(str)
	}

	return controller.ValidateString(str, "username")
}
