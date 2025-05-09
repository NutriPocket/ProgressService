// Package repository provides structs and methods to interact with the database.
package repository

import (
	"errors"
	"time"

	"github.com/NutriPocket/UserService/database"
	"github.com/NutriPocket/UserService/model"
	"github.com/go-sql-driver/mysql"
)

// IJWTRepository is an interface that contains the methods that will implement a repository struct that interact with the jwt_blacklist table.
type IJWTRepository interface {
	// Blacklist adds a JWT signature to the blacklist table.
	// signature is the JWT signature to blacklist.
	// expiresAt is the time when the JWT will expire.
	// It returns an error if the operation fails.
	Blacklist(signature string, expiresAt time.Time) error
	// IsBlacklisted checks if a JWT signature is blacklisted.
	// signature is the JWT signature to check.
	// It returns true if the JWT signature is blacklisted, false otherwise.
	// It returns an error if the operation fails.
	IsBlacklisted(signature string) (bool, error)
}

type JWTRepository struct {
	db IDatabase
}

func NewJWTRepository(db IDatabase) (*JWTRepository, error) {
	var err error

	if db == nil {
		db, err = database.GetPoolConnection()
		if err != nil {
			log.Errorf("Failed to connect to database")
			return nil, err
		}
	}

	return &JWTRepository{
		db: db,
	}, nil
}

func (r *JWTRepository) Blacklist(signature string, expiresAt time.Time) error {
	res := r.db.Exec(`
		INSERT INTO jwt_blacklist (signature, expires_at)
		VALUES (?, ?);
	`, signature, expiresAt)

	if res.Error != nil {
		if errors.Is(res.Error, &mysql.MySQLError{Number: 1062}) {
			return &model.EntityAlreadyExistsError{
				Title:  "Token no longer used",
				Detail: "The provided token is no longer in use",
			}
		}

		return res.Error
	}

	return nil
}

func (r *JWTRepository) IsBlacklisted(signature string) (bool, error) {
	var blacklistedJWT struct{ Signature string }

	res := r.db.Raw(`
		SELECT signature
		FROM jwt_blacklist
		WHERE signature = ?
	`, signature).Scan(&blacklistedJWT)

	if res.Error != nil {
		return true, res.Error
	}

	return blacklistedJWT == struct{ Signature string }{Signature: signature}, nil
}
