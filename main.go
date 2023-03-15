package main

import (
	"log"

	"github.com/abdulkarimogaji/blognado/api"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// init env variables
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// connect db
	err = db.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}

	// run server
	err = api.RunServer()
	if err != nil {
		log.Fatalf("failed to run server %v", err)
	}
}
