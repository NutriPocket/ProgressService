package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
)

func CreateUser(userData *model.BaseUser) model.User {
	hashPasswordBytes := sha256.Sum256([]byte(userData.Password))
	userData.Password = hex.EncodeToString(hashPasswordBytes[:])

	return repository.CreateUser(userData)
}

func GetAllUsers() []model.User {
	return repository.GetAllUsers()
}
