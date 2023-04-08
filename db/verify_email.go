package db

import (
	"fmt"
	"time"
)

type CreateVerifyEmailRequest struct {
	UserId     int
	SecretCode string
	Email      string
}
type VerifyEmailRequest struct {
	SecretCode string `form:"secret_code" binding:"required,min=32"`
	Email      string `form:"email" binding:"required,email"`
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

func (d *DBStruct) VerifyEmail(body VerifyEmailRequest) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// get verify_email
	row := tx.QueryRow("SELECT id, user_id, is_used, expired_at from verify_emails WHERE email = ? AND secret_code = ?", body.Email, body.SecretCode)

	var verifyEmail VerifyEmail
	err = row.Scan(&verifyEmail.Id, &verifyEmail.UserId, &verifyEmail.IsUsed, &verifyEmail.ExpiredAt)
	if err != nil {
		return err
	}

	// check if expired
	if time.Now().After(verifyEmail.ExpiredAt) {
		return fmt.Errorf("code has expired")
	}

	// set is_used to true
	stmt, err := tx.Prepare("UPDATE verify_email SET is_used = true WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(verifyEmail.Id)
	if err != nil {
		return err
	}

	// set is_email_verified to true
	stmt, err = tx.Prepare("UPDATE users SET is_email_verified = true WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(verifyEmail.UserId)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
