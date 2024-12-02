package user

import (
	"database/sql"
	"errors"
	"fmt"
	"oravue_backend/db"
	"oravue_backend/internal/config"
	"time"
)

type UserRepository interface {
	GetOtpRepository(phoneNumber string, otp string, expiry time.Time) (string, error)
	VerifyOtpRepository(phoneNumber string, otp string, config *config.Config) (string, error)
}

type UserRepoStruct struct {
	Db *db.Postgresql
}

func (u *UserRepoStruct) GetOtpRepository(phoneNumber string, otp string, expiry time.Time) (string, error) {
	query := `
		INSERT INTO users (phone_number, otp, expiry)
		VALUES ($1, $2, $3)
		ON CONFLICT (phone_number)
		DO UPDATE SET 
			otp = EXCLUDED.otp,
			updated_at = CURRENT_TIMESTAMP,
			expiry = EXCLUDED.expiry;
	`

	_, err := u.Db.Db.Exec(query, phoneNumber, otp, expiry)
	if err != nil {
		return " ", fmt.Errorf("failed to create user: %w", err)
	}
	return otp, nil
}

func (u *UserRepoStruct) VerifyOtpRepository(phoneNumber string, otp string, config *config.Config) (string, error) {

	var userID string
	var err error

	if config.Env == "development" {
		query := `
		SELECT id
		FROM users
		WHERE phone_number = $1 
		  `
		err = u.Db.Db.QueryRow(query, phoneNumber).Scan(&userID)
	} else {
		query := `
		SELECT id
		FROM users
		WHERE phone_number = $1 
		  AND otp = $2 
-- 		  AND expiry >= CURRENT_TIMESTAMP 
		  `
		err = u.Db.Db.QueryRow(query, phoneNumber, otp).Scan(&userID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("otp is in valid: %w", err) // OTP is invalid or expired
		}
		return "", fmt.Errorf("failed to query OTP: %w", err)
	}

	return userID, nil
}
