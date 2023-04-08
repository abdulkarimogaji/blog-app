package main

import (
	"log"

	"github.com/abdulkarimogaji/blognado/api"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/worker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
)

func main() {
	// init env variables
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// connect db
	dbService, err := db.NewDBService()
	if err != nil {
		log.Fatalf("failed to connect to db %v", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.AppConfig.REDIS_ADDRESS,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, dbService)

	// run server
	server := api.NewServer(dbService, taskDistributor)
	err = server.Start()
	if err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, dbService db.DBService) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, dbService)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("Failed to start task processor ", err)
	}
}
