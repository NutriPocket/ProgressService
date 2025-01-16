package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
)

type UserService struct {}

func (service *UserService) CreateUser(userData *model.BaseUser) model.User {
	hashPasswordBytes := sha256.Sum256([]byte(userData.Password))
	userData.Password = hex.EncodeToString(hashPasswordBytes[:])

	return repository.CreateUser(userData)
}

func (service *UserService) GetAllUsers() []model.User {
	return repository.GetAllUsers()
}

func (service *UserService) GetUser(username string) (model.User, error) {
	user := repository.GetUser(username)

	if user == (model.User{}) {
		return user, &model.NotFoundError{Title: "User not found", Detail: "The user with the username " + username + " was not found"}
	}

	return user, nil
}
