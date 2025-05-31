// Package database provides connection functions to the database.
package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("log")
var db *sql.DB

// ConnectDB connects to the database.
// If it fails to connect to the database, it will try again 5 times. If it fails all 5 times, it will panic.
// If it connects to the database, it will print a message to the console and assign the DB variable to the connection.
func ConnectDB() {
	if db != nil {
		return
	}

	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	if db_host == "" {
		db_host = "0.0.0.0"
	}

	if db_port == "" {
		db_port = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", db_user, db_password, db_host, db_port, db_name)

	log.Infof("DSN: %s\n", dsn)

	var try uint
	var err error

	for try < 5 {
		db, err = sql.Open("mysql", dsn)

		if err != nil {
			log.Infof("Failed to connect to database, trying again. Try number: %d\n. Err: %v", try, err)
			time.Sleep(2 * time.Second)
			try++
			continue
		}

		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)

		log.Info("Connected to database")
		return
	}

	log.Panicf("Failed to connect to database, %s", err)
}

func GetPoolConnection() (*gorm.DB, error) {
	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{},
	)

	if err != nil {
		log.Errorf("Failed to initialize gorm: %v", err)
		return nil, err
	}

	return gormDB, nil
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Errorf("Failed to close database connection: %s\n", err)
			return
		}
	}

	log.Infof("Database connection closed")
}
