// Package repository provides structs and methods to interact with the database.
package repository

import (
	"errors"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/go-sql-driver/mysql"
)

// IUserRepository is an interface that contains the methods that will implement a repository struct that interact with the users table.
type IUserRepository interface {
	// CreateUser creates a new user in the database.
	// userData is the user data to create.
	// It returns the created user and an error if the operation fails.
	CreateUser(userData *model.BaseUser) (model.User, error)
	// GetUser gets a user from the database.
	// username is the username of the user to get.
	// It returns the user and an error if the operation fails.
	GetUser(username string) (model.User, error)
	// GetUserWithPassword gets a user with the password from the database.
	// emailOrUsername is the email or username of the user to get.
	// It returns the user and an error if the operation fails.
	GetUserWithPassword(emailOrUsername string) (model.BaseUser, error)
	// GetAllUsers gets all the users from the database.
	// It returns all the users and an error if the operation fails.
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
