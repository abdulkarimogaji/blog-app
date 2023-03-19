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
	SignUp(body SignUpRequest) (int64, error)
	GetUserByEmail(email string) (User, error)
}

func (db *DBStruct) PingDB() error {
	return db.DB.Ping()
}

func NewDBService() (DBService, error) {
	conn, err := sql.Open("mysql", config.AppConfig.DB_URI)
	return &DBStruct{DB: conn}, err
}
