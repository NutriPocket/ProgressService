package repository

import (
	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/model"
)

func CreateUser(userData *model.BaseUser) model.User {
	var user model.User

	database.DB.Exec(`
			INSERT INTO users (username, email, password) 
			VALUES (?, ?, ?);
		`,
		userData.Username, userData.Email, userData.Password,
	)

	database.DB.Raw("SELECT username, email FROM users WHERE username = ?", userData.Username).Scan(&user)

	return user
}

func GetUser() {

}

func GetAllUsers() {

}
