package controller

import (
	"strings"
	"testing"
)

func TestValidateString(t *testing.T) {
	t.Run("A valid string", func(t *testing.T) {
		controller := UserController{}
		str := "This is a valid string!"

		err := controller.ValidateString(str, "valid string")

		if err != nil {
			t.Errorf("The string '%s' is invalid, what?", str)
		}
	})

	t.Run("An empty string is invalid", func(t *testing.T) {
		controller := UserController{}
		str := ""

		err := controller.ValidateString(str, "an empty string")

		if err == nil {
			t.Errorf("The string '%s' is valid, what?", str)
		}
	})

	t.Run("A string longer than 100 characters is invalid", func(t *testing.T) {
		controller := UserController{}
		str := strings.Repeat("char", 30)

		err := controller.ValidateString(str, "a long string")

		if err == nil {
			t.Errorf("The string '%s' is valid, what?", str)
		}
	})
}

func TestValidateEmail(t *testing.T) {
	t.Run("A valid email", func(t *testing.T) {
		controller := UserController{}
		email := "test@test.com"

		err := controller.ValidateEmail(email)

		if err != nil {
			t.Errorf("The email '%s' is invalid, what?", email)
		}
	})

	t.Run("An empty email is invalid", func(t *testing.T) {
		controller := UserController{}
		email := ""

		err := controller.ValidateEmail(email)

		if err == nil {
			t.Errorf("The email '%s' is valid, what?", email)
		}
	})

	t.Run("A email longer than 100 characters is invalid", func(t *testing.T) {
		controller := UserController{}
		email := strings.Repeat("test@test.com", 20)

		err := controller.ValidateEmail(email)

		if err == nil {
			t.Errorf("The email '%s' is valid, what?", email)
		}
	})

	t.Run("A email without @ is invalid", func(t *testing.T) {
		controller := UserController{}
		email := "testtest.com"

		err := controller.ValidateEmail(email)

		if err == nil {
			t.Errorf("The email '%s' is valid, what?", email)
		}
	})
}
