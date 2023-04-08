package db

import (
	"time"
)

type User struct {
	Id              int       `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
	Password        string    `json:"password,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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

type Blog struct {
	Id        int        `json:"id"`
	AuthorId  int        `json:"author_id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Excerpt   string     `json:"excerpt"`
	Thumbnail string     `json:"thumbnail"`
	Body      string     `json:"body,omitempty"`
	PostedAt  time.Time  `json:"posted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Author    BlogAuthor `json:"author,omitempty"`
}

type BlogAuthor struct {
	Id        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email,omitempty"`
	Photo     *string `json:"photo,omitempty"`
	Socials   *string `json:"socials,omitempty"`
}

type Comment struct {
	Id        int         `json:"id"`
	UserId    int         `json:"user_id"`
	BlogId    int         `json:"blog_id"`
	Message   string      `json:"message"`
	Thread    string      `json:"thread"`
	PostedAt  time.Time   `json:"posted_at"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	User      CommentUser `json:"user,omitempty"`
}

type CommentUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string `json:"photo"`
}

type VerifyEmail struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Email      string    `json:"email"`
	IsUsed     bool      `json:"is_used"`
	SecretCode string    `json:"secret_code"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
