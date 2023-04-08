package worker

import (
	"context"
	"log"

	"github.com/abdulkarimogaji/blognado/api/lambda"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server    *asynq.Server
	dbService db.DBService
	mailer    lambda.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, dbService db.DBService, mailer lambda.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Println("ERROR LOG: ", task.Type(), err)
		}),
	})
	return &RedisTaskProcessor{
		server:    server,
		dbService: dbService,
		mailer:    mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
