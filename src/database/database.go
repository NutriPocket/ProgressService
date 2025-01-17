package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(0.0.0.0:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", db_user, db_password, db_name)

	var try uint
	var err error

	for try < 5 {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Printf("Failed to connect to database, trying again. Try number: %d\n", try)
			time.Sleep(2 * time.Second)
			try++
			continue
		}

		log.Println("Connected to database")
		return
	}

	panic(fmt.Sprintf("Failed to connect to database, %s", err))
}
