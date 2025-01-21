package model

import "fmt"

type ValidationError struct {
	Detail string
	Title  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s, %s", e.Title, e.Detail)
}

type AuthenticationError struct {
	Detail string
	Title  string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%s, %s", e.Title, e.Detail)
}

type NotFoundError struct {
	Detail string
	Title  string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s, %s", e.Title, e.Detail)
}

type EntityAlreadyExistsError struct {
	Detail string
	Title  string
}

func (e *EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s, %s", e.Title, e.Detail)
}
