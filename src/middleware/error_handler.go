package middleware

import (
	"net/http"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/gin-gonic/gin"
)

type errorRfc9457 struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()

		if err != nil {
			var status int
			var detail string
			var title string

			switch e := err.Err.(type) {
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
			default:
				status = http.StatusInternalServerError
				detail = "An unknown error has occurred"
				title = "Internal Server Error"
			}

			rfcError := errorRfc9457{
				Type:     "about:blank",
				Title:    title,
				Status:   status,
				Detail:   detail,
				Instance: c.Request.URL.Path,
			}

			c.JSON(rfcError.Status, rfcError)

			c.Abort()
		}
	}
}
