package middleware

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/MaxiOtero6/go-auth-rest/model"
)

func TestParseError(t *testing.T) {
	t.Run("An unknown error is parsed as an internal server error with status code 500", func(t *testing.T) {
		urlPath := "/"

		expected := errorRfc9457{
			Title:    "Internal Server Error",
			Detail:   "An unknown error has occurred",
			Status:   http.StatusInternalServerError,
			Type:     "about:blank",
			Instance: "/",
		}

		err := errors.New("Unknown error :)")

		result := parseError(err, urlPath)

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("The parsed error isn't equal to the expected one")
		}
	})

	t.Run("A validation error is parsed with status code 400", func(t *testing.T) {
		urlPath := "/"

		detail := "The specified date has an unknown format, try: YYYY-MM-DD"
		title := "Unknown date format"

		expected := errorRfc9457{
			Title:    title,
			Detail:   detail,
			Status:   http.StatusBadRequest,
			Type:     "about:blank",
			Instance: "/",
		}

		err := &model.ValidationError{
			Title:  title,
			Detail: detail,
		}

		result := parseError(err, urlPath)

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("The parsed error isn't equal to the expected one")
		}
	})

	t.Run("An authentication error is parsed with status code 401", func(t *testing.T) {
		urlPath := "/"

		detail := "The specified token has an unknown format, try: Bearer <token>"
		title := "Unknown token format"

		expected := errorRfc9457{
			Title:    title,
			Detail:   detail,
			Status:   http.StatusUnauthorized,
			Type:     "about:blank",
			Instance: "/",
		}

		err := &model.AuthenticationError{
			Title:  title,
			Detail: detail,
		}

		result := parseError(err, urlPath)

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("The parsed error isn't equal to the expected one")
		}
	})

	t.Run("A not found error is parsed with status code 404", func(t *testing.T) {
		urlPath := "/"

		detail := "The specified user is missing, please try another username"
		title := "User not found"

		expected := errorRfc9457{
			Title:    title,
			Detail:   detail,
			Status:   http.StatusNotFound,
			Type:     "about:blank",
			Instance: "/",
		}

		err := &model.NotFoundError{
			Title:  title,
			Detail: detail,
		}

		result := parseError(err, urlPath)

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("The parsed error isn't equal to the expected one")
		}
	})

	t.Run("An entity already exists error is parsed with status code 409", func(t *testing.T) {
		urlPath := "/"

		detail := "The specified user already exists, please try another username"
		title := "User already exists"

		expected := errorRfc9457{
			Title:    title,
			Detail:   detail,
			Status:   http.StatusConflict,
			Type:     "about:blank",
			Instance: "/",
		}

		err := &model.EntityAlreadyExistsError{
			Title:  title,
			Detail: detail,
		}

		result := parseError(err, urlPath)

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("The parsed error isn't equal to the expected one")
		}
	})
}
