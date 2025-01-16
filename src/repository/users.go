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

func GetUser(username string) model.User {
	var user model.User

	database.DB.Raw("SELECT username, email FROM users WHERE username = ?", username).Scan(&user)

	return user
}

func GetUserWithPassword(emailOrUsername string) model.BaseUser {
	var user model.BaseUser

	database.DB.Raw("SELECT username, email, password FROM users WHERE username = ? OR email = ?", emailOrUsername, emailOrUsername).Scan(&user)

	return user
}

func GetAllUsers() []model.User {
	var users []model.User

	database.DB.Raw("SELECT username, email FROM users").Scan(&users)

	return users
}
