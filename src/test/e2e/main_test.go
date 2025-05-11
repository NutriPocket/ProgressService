package e2e_test

import (
	"os"
	"testing"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/test"
	"github.com/NutriPocket/ProgressService/utils"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")
var router *gin.Engine
var bearerToken string

var testUser = model.User{
	ID: "1", Username: "test", Email: "test@test.com",
}

func TestMain(m *testing.M) {
	test.Setup("e2e")
	gin.SetMode(gin.TestMode)

	router = utils.SetupRouter()
	bearerToken = test.GetBearerToken(&testUser)

	log.Infof("Running e2e tests with token: %s\n", bearerToken)

	code := m.Run()
	test.TearDown("e2e")
	os.Exit(code)
}
