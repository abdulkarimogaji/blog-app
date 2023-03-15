package db

import (
	"database/sql"

	"github.com/abdulkarimogaji/blognado/config"
)

type DBService struct {
}

var DbConn *sql.DB
var DbService DBService

func ConnectDB() error {
	var err error
	DbConn, err = sql.Open("mysql", config.AppConfig.DB_URI)
	return err
}
