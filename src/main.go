package main

import (
	"os"

	"github.com/NutriPocket/ProgressService/database"
	"github.com/NutriPocket/ProgressService/utils"
	"github.com/joho/godotenv"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

func loadEnv() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "../.env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// InitLogger Receives the log level to be set in go-logging as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	baseBackend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05} %{level:.5s}     %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(baseBackend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	logLevelCode, err := logging.LogLevel(logLevel)
	if err != nil {
		return err
	}
	backendLeveled.SetLevel(logLevelCode, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
	log.Infof("Log level set to %s", logLevel)
	return nil
}

func main() {
	loadEnv()

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}

	InitLogger(logLevel)

	database.ConnectDB()
	defer database.Close()

	router := utils.SetupRouter()

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	addr := host + ":" + port

	log.Infof("Starting server on %s", addr)
	router.Run(addr)
}
