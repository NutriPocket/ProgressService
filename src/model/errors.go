// Package model contains the structs types that will be used in the application.
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

type ConflictError struct {
	Detail string
	Title  string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s, %s", e.Title, e.Detail)
}
