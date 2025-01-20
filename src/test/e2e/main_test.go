package e2e_test

import (
	"os"
	"testing"

	"github.com/MaxiOtero6/go-auth-rest/test"
	"github.com/MaxiOtero6/go-auth-rest/utils"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	test.Setup("e2e")
	gin.SetMode(gin.TestMode)
	router = utils.SetupRouter()
	code := m.Run()
	test.TearDown("e2e")
	os.Exit(code)
}
