package repository

import (
	"errors"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/go-sql-driver/mysql"
)

type IUserRepository interface {
	CreateUser(userData *model.BaseUser) (model.User, error)
	GetUser(username string) (model.User, error)
	GetUserWithPassword(emailOrUsername string) (model.BaseUser, error)
	GetAllUsers() ([]model.User, error)
}

type UserRepository struct{}

func (repository *UserRepository) CreateUser(userData *model.BaseUser) (model.User, error) {
	var user model.User

	res := database.DB.Exec(`
			INSERT INTO users (username, email, password) 
			VALUES (?, ?, ?);
		`,
		userData.Username, userData.Email, userData.Password,
	)

	if res.Error != nil {
		if errors.Is(res.Error, &mysql.MySQLError{Number: 1062}) {
			return model.User{}, &model.EntityAlreadyExistsError{
				Title:  "Username or email already in use",
				Detail: "The provided username or email are already in use, try something else",
			}
		}

		return model.User{}, res.Error
	}

	res = database.DB.Raw("SELECT username, email FROM users WHERE username = ?", userData.Username).Scan(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}

func (repository *UserRepository) GetUser(username string) (model.User, error) {
	var user model.User

	res := database.DB.Raw("SELECT username, email FROM users WHERE username = ?", username).Scan(&user)

	if res.Error != nil {
		return model.User{}, res.Error
	}

	return user, nil
}

func (repository *UserRepository) GetUserWithPassword(emailOrUsername string) (model.BaseUser, error) {
	var user model.BaseUser

	res := database.DB.Raw("SELECT username, email, password FROM users WHERE username = ? OR email = ?", emailOrUsername, emailOrUsername).Scan(&user)

	if res.Error != nil {
		return model.BaseUser{}, res.Error
	}

	return user, nil
}

func (repository *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User

	res := database.DB.Raw("SELECT username, email FROM users ORDER BY created_at DESC").Scan(&users)

	if res.Error != nil {
		return []model.User{}, res.Error
	}

	return users, nil
}
