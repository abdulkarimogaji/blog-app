package db

import "time"

type User struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Photo       string    `json:"photo"`
	DateOfBirth time.Time `json:"date_of_birth"`
	About       string    `json:"about"`
	Settings    string    `json:"settings"`
	Socials     string    `json:"socials"`
}
