// Package database provides connection functions to the database.
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

// ConnectDB connects to the database.
// If it fails to connect to the database, it will try again 5 times. If it fails all 5 times, it will panic.
// If it connects to the database, it will print a message to the console and assign the DB variable to the connection.
func ConnectDB() {
	if DB != nil {
		return
	}

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	db_host := os.Getenv("DB_HOST")

	if db_host == "" {
		db_host = "localhost"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", db_user, db_password, db_host, db_name)

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
