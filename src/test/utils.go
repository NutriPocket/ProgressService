package test

import (
	"log"
	"os"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/joho/godotenv"
)

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

	if err := database.DB.Exec("CREATE DATABASE IF NOT EXISTS test").Error; err != nil {
		log.Fatal(err)
	}

	if err := database.DB.Exec("USE test").Error; err != nil {
		log.Fatal(err)
	}

	if err := database.DB.Exec(`
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

	if err := database.DB.Exec(`
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
	if err := database.DB.Exec(`
		DELETE FROM users
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func ClearBlacklist() {
	if err := database.DB.Exec(`
		DELETE FROM jwt_blacklist
	`).Error; err != nil {
		log.Fatal(err)
	}
}

func Setup(testType string) {
	log.Printf("Setup %s tests!\n", testType)
	loadEnv()
	log.Println(".env.test loaded")
	setupDB()
}

func TearDown(testType string) {
	log.Printf("Tear down %s tests!\n", testType)
}
