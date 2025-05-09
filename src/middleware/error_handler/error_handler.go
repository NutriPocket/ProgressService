// Package middleware provides custom middlewares for the API
package middleware

import (
	"net/http"

	"github.com/NutriPocket/UserService/model"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// errorRfc9457 is a struct that will be used to return errors in the RFC 9457 format
type errorRfc9457 struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// parseError parses an error and returns an error in the RFC 9457 format
// err is the error to parse
// urlPath is the URL path of the request
// It returns an error in the RFC 9457 format
func parseError(err error, urlPath string) errorRfc9457 {
	var status int
	var detail string
	var title string

	switch e := err.(type) {
	case *model.ValidationError:
		status = http.StatusBadRequest
		detail = e.Detail
		title = e.Title
	case *model.AuthenticationError:
		status = http.StatusUnauthorized
		detail = e.Detail
		title = e.Title
	case *model.NotFoundError:
		status = http.StatusNotFound
		detail = e.Detail
		title = e.Title
	case *model.EntityAlreadyExistsError:
		status = http.StatusConflict
		detail = e.Detail
		title = e.Title
	default:
		status = http.StatusInternalServerError
		detail = "An unknown error has occurred"
		title = "Internal Server Error"
	}

	return errorRfc9457{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: urlPath,
	}
}

// ErrorHandler is a middleware that handles errors and returns them in the RFC 9457 format
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()

		if err != nil {
			rfcError := parseError(err.Err, c.Request.URL.Path)
			log.Errorf("Error: %s", rfcError)

			c.JSON(rfcError.Status, rfcError)

			c.Abort()
		}
	}
}
