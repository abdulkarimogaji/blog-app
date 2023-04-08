package db

import (
	"time"
)

type SignUpRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	City        *string `json:"city"`
	Country     *string `json:"country"`
	Photo       *string `json:"photo" binding:"omitempty,url"`
	DateOfBirth *string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	About       *string `json:"about"`
	Settings    *string `json:"settings" binding:"omitempty,json"`
	Socials     *string `json:"socials" binding:"omitempty,json"`
}

func (d *DBStruct) SignUp(body SignUpRequest) (int64, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO users (first_name, last_name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return 0, err
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	result, err := stmt.Exec(body.FirstName, body.LastName, body.Email, body.Password, createdAt, updatedAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// insert to profile table
	stmt, err = tx.Prepare("INSERT INTO profile (user_id, date_of_birth, about, photo, city, country, settings, socials, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")

	if err != nil {
		return id, err
	}

	result, err = stmt.Exec(id, body.DateOfBirth, body.About, body.Photo, body.City, body.Country, body.Settings, body.Socials, createdAt, updatedAt)

	if err != nil {
		return id, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return id, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}
