package controller

import (
	"regexp"

	"github.com/MaxiOtero6/go-auth-rest/model"
)

func ValidateString(str string, field string) error {
	if str == "" {
		return &model.ValidationError{Detail: "The " + field + " field is required", Title: "Empty " + field + " field"}
	}

	if len(str) > 100 {
		return &model.ValidationError{Detail: "The " + field + " field must be less than 100 characters", Title: "Invalid " + field + " field"}
	}

	return nil
}

func ValidateEmail(email string) error {
	var err error

	if err = ValidateString(email, "email"); err != nil {
		return err
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`)

	if !emailRegex.MatchString(email) {
		return &model.ValidationError{Detail: "The email field must be a valid email address", Title: "Invalid email field"}
	}

	return nil
}
