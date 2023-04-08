package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/util"
	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Println("INFO: ", task.Type(), info.Queue, info.MaxRetry)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.dbService.GetUserByEmail(payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	secretCode := util.RandomString(32)

	_, err = processor.dbService.CreateVerifyEmail(db.CreateVerifyEmailRequest{
		UserId:     user.Id,
		Email:      user.Email,
		SecretCode: secretCode,
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}

	err = processor.mailer.SendEmail("Verify email", fmt.Sprintf(`Please verify your email using this <a href="http://localhost:4000/verify-email?code=%s">link</a>`, secretCode), []string{user.Email, "abdulkarimogaji001@gmail.com"}, nil, nil, nil)

	if err != nil {
		return fmt.Errorf("failed to send verify email %w", err)
	}
	log.Println("INFO LOG: ", task.Type(), user.FirstName, user.LastName)

	return nil
}
