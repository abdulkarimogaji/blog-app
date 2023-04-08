package db

import (
	"time"
)

type CreateVerifyEmailRequest struct {
	UserId     int
	SecretCode string
	Email      string
}

func (d *DBStruct) CreateVerifyEmail(body CreateVerifyEmailRequest) (int, error) {

	stmt, err := d.DB.Prepare("INSERT INTO verify_emails (user_id, email, secret_code, created_at, expired_at) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return 0, err
	}

	createdAt := time.Now()
	expiredAt := time.Now().Add(15 * time.Minute)

	result, err := stmt.Exec(body.UserId, body.Email, body.SecretCode, createdAt, expiredAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
