package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
)

type UserService struct{}

func (service *UserService) EncodePassword(password string) string {
	hashPasswordBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashPasswordBytes[:])

}

func (service *UserService) CreateUser(userData *model.BaseUser) model.User {
	userData.Password = service.EncodePassword(userData.Password)

	return repository.CreateUser(userData)
}

func (service *UserService) Login(userData *model.LoginUser) (model.User, error) {
	savedUser := repository.GetUserWithPassword(userData.EmailOrUsername)

	if savedUser == (model.BaseUser{}) {
		return model.User{}, &model.AuthenticationError{
			Title:  "Credentials don't match",
			Detail: "User identification or password are wrong, please try again",
		}
	}

	userData.Password = service.EncodePassword(userData.Password)

	if userData.Password != savedUser.Password {
		return model.User{}, &model.AuthenticationError{
			Title:  "Credentials don't match",
			Detail: "User identification or password are wrong, please try again",
		}
	}

	return model.User{Username: savedUser.Username, Email: savedUser.Email}, nil
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
