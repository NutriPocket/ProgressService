package main

import (
	"log"
	"os"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/middleware"
	"github.com/MaxiOtero6/go-auth-rest/routes"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.Use(middleware.ErrorHandler())
	router.Use(middleware.AuthMiddleware())
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	addr := host + ":" + port

	router.Run(addr)
}
