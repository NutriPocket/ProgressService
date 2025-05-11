package middleware

import (
	"reflect"
	"testing"

	"github.com/NutriPocket/ProgressService/model"
)

func TestGetRootPath(t *testing.T) {
	t.Run("An empty url should return an empty string", func(t *testing.T) {
		url := ""
		expected := ""

		result := getRootPath(url)

		if expected != result {
			t.Errorf("getRootPath should be an '%s'", expected)
		}
	})

	t.Run("A root url should return an empty string", func(t *testing.T) {
		url := "/"
		expected := ""

		result := getRootPath(url)

		if expected != result {
			t.Errorf("getRootPath should be an '%s'", expected)
		}
	})

	t.Run("A large url should return the first resource name", func(t *testing.T) {
		url := "/users/testUser/followers"
		expected := "users"

		result := getRootPath(url)

		if expected != result {
			t.Errorf("getRootPath should be '%s'", expected)
		}
	})

	t.Run("A simple resource url should return the resource name", func(t *testing.T) {
		url := "/users/"
		expected := "users"

		result := getRootPath(url)

		if expected != result {
			t.Errorf("getRootPath should be '%s'", expected)
		}
	})
}

func TestGetToken(t *testing.T) {
	t.Run("It should return the parsed token", func(t *testing.T) {
		authHeader := "Bearer this-is-a-token"
		expected := "this-is-a-token"

		token, err := getToken(authHeader)

		if err != nil {
			t.Error(err)
		}

		if expected != token {
			t.Errorf("The parsed token should be %s", expected)
		}
	})

	t.Run("It should return an authentication error if the authHeader is an empty string", func(t *testing.T) {
		authHeader := ""
		expected := &model.AuthenticationError{
			Title:  "Unauthorized user",
			Detail: "The user isn't authorized because no Authorization header is provided",
		}

		token, err := getToken(authHeader)

		if token != "" {
			t.Error("The token should be an empty string")
		}

		if !reflect.DeepEqual(expected, err) {
			t.Errorf("It should return the following error: %s", expected)
		}
	})

	t.Run("It should return an authentication error if the authHeader has an invalid format (no spaces)", func(t *testing.T) {
		authHeader := "Bearerthis-is-a-token"
		expected := &model.AuthenticationError{
			Title:  "Invalid authorization",
			Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
		}

		token, err := getToken(authHeader)

		if token != "" {
			t.Error("The token should be an empty string")
		}

		if !reflect.DeepEqual(expected, err) {
			t.Errorf("It should return the following error: %s", expected)
		}
	})

	t.Run("It should return an authentication error if the authHeader has an invalid format (more than one space)", func(t *testing.T) {
		authHeader := "Bearer  this-is-a-token"
		expected := &model.AuthenticationError{
			Title:  "Invalid authorization",
			Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
		}

		token, err := getToken(authHeader)

		if token != "" {
			t.Error("The token should be an empty string")
		}

		if !reflect.DeepEqual(expected, err) {
			t.Errorf("It should return the following error: %s", expected)
		}
	})

	t.Run("It should return an authentication error if the authHeader has an invalid format (missing Bearer)", func(t *testing.T) {
		authHeader := " this-is-a-token"
		expected := &model.AuthenticationError{
			Title:  "Invalid authorization",
			Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
		}

		token, err := getToken(authHeader)

		if token != "" {
			t.Error("The token should be an empty string")
		}

		if !reflect.DeepEqual(expected, err) {
			t.Errorf("It should return the following error: %s", expected)
		}
	})

	t.Run("It should return an authentication error if the authHeader has an invalid format (Bearer typo)", func(t *testing.T) {
		authHeader := "Biarer this-is-a-token"
		expected := &model.AuthenticationError{
			Title:  "Invalid authorization",
			Detail: "The Authorization header provided has an unknown format... Try: Bearer <token>",
		}

		token, err := getToken(authHeader)

		if token != "" {
			t.Error("The token should be an empty string")
		}

		if !reflect.DeepEqual(expected, err) {
			t.Errorf("It should return the following error: %s", expected)
		}
	})
}
