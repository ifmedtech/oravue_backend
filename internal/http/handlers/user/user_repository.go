package user

import (
	"fmt"
	"oravue_backend/db"
	"time"
)

type UserRepository interface {
	GetOtpRepository(phoneNumber string, otp string, expiry time.Time) (string, error)
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
