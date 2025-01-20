package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
)

type UserService struct {
	repository repository.IUserRepository
}

func NewUserService(userRepository *repository.IUserRepository) UserService {
	var repo repository.IUserRepository = &repository.UserRepository{}

	if userRepository != nil {
		repo = *userRepository
	}

	return UserService{repository: repo}
}

func (service *UserService) EncodePassword(password string) string {
	hashPasswordBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashPasswordBytes[:])

}

func (service *UserService) CreateUser(userData *model.BaseUser) (model.User, error) {
	userData.Password = service.EncodePassword(userData.Password)

	return service.repository.CreateUser(userData)
}

func (service *UserService) Login(userData *model.LoginUser) (model.User, error) {
	savedUser, err := service.repository.GetUserWithPassword(userData.EmailOrUsername)

	if err != nil {
		return model.User{}, err
	}

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

func (service *UserService) GetAllUsers() ([]model.User, error) {
	return service.repository.GetAllUsers()
}

func (service *UserService) GetUser(username string) (model.User, error) {
	user, err := service.repository.GetUser(username)

	if err != nil {
		return user, err
	}

	if user == (model.User{}) {
		return user, &model.NotFoundError{Title: "User not found", Detail: "The user with the username " + username + " was not found"}
	}

	return user, nil
}
