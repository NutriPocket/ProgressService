package test

import (
	"fmt"
	"os"

	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/service"
	"github.com/joho/godotenv"
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("log")
var gormDB *gorm.DB

func loadEnv() {
	if ci_test := os.Getenv("CI_TEST"); ci_test != "" {
		return
	}

	err := godotenv.Load("../../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func setupDB() {
	database.ConnectDB()
	var err error
	gormDB, err = database.GetPoolConnection()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}
}

func ClearAllData() {
	if err := gormDB.Exec(`
		DELETE FROM fixed_user_data;
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec(`
		DELETE FROM anthropometric_data;
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec(`
		DELETE FROM objective;
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec(`
		DELETE FROM user_routines;
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func Setup(testType string) {
	log.Infof("Setup %s tests!\n", testType)
	loadEnv()
	log.Info(".env.test loaded")
	setupDB()
}

func TearDown(testType string) {
	log.Infof("Tear down %s tests!\n", testType)
	database.Close()
}

func GetBearerToken(testUser *model.User) string {
	jwtService, err := service.NewJWTService()
	if err != nil {
		log.Fatalf("An error ocurred when creating the JWT service: %v\n", err)
	}

	token, err := jwtService.Sign(*testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	return fmt.Sprintf("Bearer %s", token)
}
