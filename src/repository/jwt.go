package repository

import (
	"errors"
	"time"

	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/go-sql-driver/mysql"
)

type IJWTRepository interface {
	Blacklist(signature string, expiresAt time.Time) error
	IsBlacklisted(signature string) (bool, error)
}

type JWTRepository struct{}

func (repository *JWTRepository) Blacklist(signature string, expiresAt time.Time) error {
	res := database.DB.Exec(`
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

func (repository *JWTRepository) IsBlacklisted(signature string) (bool, error) {
	var blacklistedJWT struct{ Signature string }

	res := database.DB.Raw(`
		SELECT signature
		FROM jwt_blacklist
		WHERE signature = ?
	`, signature).Scan(&blacklistedJWT)

	if res.Error != nil {
		return true, res.Error
	}

	return blacklistedJWT == struct{ Signature string }{Signature: signature}, nil
}
