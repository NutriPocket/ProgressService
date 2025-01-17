package repository

import (
	"time"

	"github.com/MaxiOtero6/go-auth-rest/database"
)

type JWTRepository struct{}

func (repository *JWTRepository) Blacklist(signature string, expiresAt time.Time) {
	database.DB.Exec(`
		INSERT INTO jwt_blacklist (signature, expires_at)
		VALUES (?, ?);
	`, signature, expiresAt)
}

func (repository *JWTRepository) IsBlacklisted(signature string) bool {
	var blacklistedJWT struct{ Signature string }

	database.DB.Raw(`
		SELECT signature
		FROM jwt_blacklist
		WHERE signature = ?
	`, signature).Scan(&blacklistedJWT)

	return blacklistedJWT == struct{ Signature string }{Signature: signature}
}
