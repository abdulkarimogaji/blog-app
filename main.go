package main

import (
	"log"

	"github.com/abdulkarimogaji/blognado/api"
	"github.com/abdulkarimogaji/blognado/api/lambda"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/worker"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-redis/redis/v8"
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

	// connect cloudinary
	cloudinaryInstance, err := cloudinary.NewFromURL(config.AppConfig.CLOUDINARY_URL)
	if err != nil {
		log.Fatalf("failed to connect to cloudinary %v", err)
	}

	opts, err := redis.ParseURL(config.AppConfig.REDIS_ADDRESS)
	if err != nil {
		log.Fatalf("invalid redis uri %v", err)
	}

	// setup async worker
	redisOpt := asynq.RedisClientOpt{
		Addr: opts.Addr,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, dbService)

	// run server
	server := api.NewServer(dbService, taskDistributor, cloudinaryInstance)
	err = server.Start()
	if err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, dbService db.DBService) {
	mailer := lambda.NewGmailSender(config.AppConfig.GMAIL_NAME, config.AppConfig.GMAIL_ADDRESS, config.AppConfig.GMAIL_PASSWORD)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, dbService, mailer)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("Failed to start task processor ", err)
	}
}
