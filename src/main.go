package main

import (
	"log"
	"os"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/utils"
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

	router := utils.SetupRouter()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	addr := host + ":" + port

	router.Run(addr)
}
