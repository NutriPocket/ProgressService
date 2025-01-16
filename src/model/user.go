package model

type BaseUser struct {
	Username string
	Email    string
	Password string
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUser struct {
	EmailOrUsername string
	Password        string
}
