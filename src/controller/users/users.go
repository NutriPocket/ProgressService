// Package controller provides the user controller struct that will be used to validate the user input received from http requests.
package controller

import (
	"regexp"
	"strings"

	"github.com/NutriPocket/ProgressService/model"
)

// UserController is a struct that will be used to validate the user input received from http requests.
type UserController struct{}

// ValidateString validates a string and returns an error if the string is empty or longer than 100 characters.
// str is he string to validate.
// field is the field name of the string.
func (controller *UserController) ValidateString(str string, field string) error {
	if str == "" {
		return &model.ValidationError{Detail: "The " + field + " field is required", Title: "Empty " + field + " field"}
	}

	if len(str) > 100 {
		return &model.ValidationError{Detail: "The " + field + " field must be less than 100 characters", Title: "Invalid " + field + " field"}
	}

	return nil
}

// ValidateEmail validates an email and returns an error if the email is not a valid email address.
// email is the email to validate.
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

// ValidateUsernameOrEmail validates a string that can be an email or a username.
// str is the string to validate.
func (controller *UserController) ValidateUsernameOrEmail(str string) error {
	if strings.Contains(str, "@") {
		return controller.ValidateEmail(str)
	}

	return controller.ValidateString(str, "username")
}
