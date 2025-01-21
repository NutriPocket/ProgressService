// Package model contains the structs types that will be used in the application.
package model

// BaseUser is a struct that contains the base user data received from the client
type BaseUser struct {
	Username string
	Email    string
	Password string
}

// User is a struct that contains the user data that will be sent to the client
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginUser is a struct that contains the user data received from the client when logging in
type LoginUser struct {
	EmailOrUsername string
	Password        string
}
