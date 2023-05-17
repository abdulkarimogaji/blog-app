package db

import (
	"context"
	"database/sql"

	"github.com/abdulkarimogaji/blognado/config"
	"github.com/google/uuid"
)

type DBStruct struct {
	DB *sql.DB
}

type DBService interface {
	PingDB() error
	SignUp(ctx context.Context, body SignUpRequest, afterCreate func(body SignUpRequest) error) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)

	CreateBlog(ctx context.Context, body CreateBlogRequest) (Blog, error)
	GetBlogBySlug(ctx context.Context, slug string) (Blog, error)
	GetBlogById(ctx context.Context, id int) (Blog, error)
	GetBlogs(ctx context.Context, filters GetBlogsFilters, params PaginationParams) ([]Blog, int, error)

	GetComments(ctx context.Context, filters GetCommentsFilters, params PaginationParams) ([]Comment, int, error)
	CreateComment(ctx context.Context, body CreateCommentRequest) (Comment, error)
	CreateVerifyEmail(ctx context.Context, body CreateVerifyEmailRequest) (int, error)
	VerifyEmail(ctx context.Context, body VerifyEmailRequest) error

	CreateSession(ctx context.Context, body CreateSessionRequest) (Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)

	DeleteRow(ctx context.Context, tableName string, id int) (int, error)
}

func (db *DBStruct) PingDB() error {
	return db.DB.Ping()
}

func NewDBService() (DBService, error) {
	conn, err := sql.Open("mysql", config.AppConfig.DB_URI)
	return &DBStruct{DB: conn}, err
}
