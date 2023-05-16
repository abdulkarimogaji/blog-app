package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/util"
	"github.com/abdulkarimogaji/blognado/worker"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
)

func signUp(dbService db.DBService, taskDistributor worker.TaskDistributor) gin.HandlerFunc {
	return func(c *gin.Context) {

		var body db.SignUpRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			validationResponse(err, c)
			return
		}

		hashedPassword, err := util.HashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		body.Password = hashedPassword
		id, err := dbService.SignUp(c, body, func(body db.SignUpRequest) error {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return taskDistributor.DistributeTaskSendVerifyEmail(context.Background(), &worker.PayloadSendVerifyEmail{Email: body.Email}, opts...)
		})
		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if ok && me.Number == MYSQL_KEY_EXISTS {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "user exists",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": id,
			"error":   false,
		})
	}
}
