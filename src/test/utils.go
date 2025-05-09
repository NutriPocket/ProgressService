package test

import (
	"os"

	"github.com/NutriPocket/UserService/database"
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

	if err := gormDB.Exec("CREATE DATABASE IF NOT EXISTS test").Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec("USE test").Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6)
		)
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := gormDB.Exec(`
		CREATE TABLE IF NOT EXISTS jwt_blacklist (
			signature VARCHAR(100) PRIMARY KEY,
			expires_at TIMESTAMP NOT NULL,
			INDEX idx_expires_at (expires_at)
		)
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func ClearUsers() {
	if err := gormDB.Exec(`
		DELETE FROM users
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func ClearBlacklist() {
	if err := gormDB.Exec(`
		DELETE FROM jwt_blacklist
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
