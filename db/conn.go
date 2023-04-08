package db

import (
	"database/sql"

	"github.com/abdulkarimogaji/blognado/config"
)

type DBStruct struct {
	DB *sql.DB
}

type DBService interface {
	PingDB() error
	SignUp(body SignUpRequest, afterCreate func(body SignUpRequest) error) (int64, error)
	GetUserByEmail(email string) (User, error)

	CreateBlog(body CreateBlogRequest) (Blog, error)
	GetBlogBySlug(slug string) (Blog, error)
	GetBlogById(id int) (Blog, error)
	GetBlogs(filters GetBlogsFilters, params PaginationParams) ([]Blog, int, error)

	GetComments(filters GetCommentsFilters, params PaginationParams) ([]Comment, int, error)
	CreateComment(body CreateCommentRequest) (Comment, error)
	CreateVerifyEmail(body CreateVerifyEmailRequest) (int, error)

	DeleteRow(tableName string, id int) (int, error)
}

func (db *DBStruct) PingDB() error {
	return db.DB.Ping()
}

func NewDBService() (DBService, error) {
	conn, err := sql.Open("mysql", config.AppConfig.DB_URI)
	return &DBStruct{DB: conn}, err
}
