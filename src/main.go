package main

import (
	"log"
	"os"

	"github.com/MaxiOtero6/go-auth-rest/database"
	middlewareAuth "github.com/MaxiOtero6/go-auth-rest/middleware/auth_middleware"
	middlewareErr "github.com/MaxiOtero6/go-auth-rest/middleware/error_handler"
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

	router.Use(middlewareErr.ErrorHandler())
	router.Use(middlewareAuth.AuthMiddleware())
	routes.AuthRoutes(router)
	routes.UsersRoutes(router)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	addr := host + ":" + port

	router.Run(addr)
}
