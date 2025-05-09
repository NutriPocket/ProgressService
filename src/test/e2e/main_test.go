package e2e_test

import (
	"os"
	"testing"

	"github.com/NutriPocket/UserService/test"
	"github.com/NutriPocket/UserService/utils"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")
var router *gin.Engine

func TestMain(m *testing.M) {
	test.Setup("e2e")
	gin.SetMode(gin.TestMode)

	router = utils.SetupRouter()

	code := m.Run()
	test.TearDown("e2e")
	os.Exit(code)
}
