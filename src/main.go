package main

import (
	"log"

	"github.com/MaxiOtero6/go-auth-rest/database"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	loadEnv()

	database.ConnectDB()
}
