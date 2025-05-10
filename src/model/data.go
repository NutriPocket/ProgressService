// Package model contains the structs types that will be used in the application.
package model

// User is a struct that contains the user data that will be received in the JWT token
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type BaseFixedUserData struct {
	UserID   string `json:"user_id"`
	Height   uint   `json:"height" binding:"required"`
	Birthday string `json:"birthday" binding:"required"`
}

type FixedUserData struct {
	UserID string `json:"user_id"`
	Height uint   `json:"height"`
	Age    uint   `json:"age"`
}

type AnthropometricData struct {
	UserID     string  `json:"user_id"`
	Weight     float32 `json:"weight" binding:"required"`
	MuscleMass float32 `json:"muscle_mass"`
	FatMass    float32 `json:"fat_mass"`
	BoneMass   float32 `json:"bone_mass"`
	CreatedAt  string  `json:"created_at"`
}
