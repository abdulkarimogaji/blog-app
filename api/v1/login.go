package v1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/util"
	"github.com/abdulkarimogaji/blognado/worker"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
)

func login(dbService db.DBService, taskDistributor worker.TaskDistributor) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body db.LoginRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)

		// request validation
		if err != nil {
			validationResponse(err, c)
			return
		}

		// get user
		user, err := dbService.GetUserByEmail(c, body.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "user does not exist",
					"error":   err.Error(),
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		//check password
		ok := util.VerifyPassword(body.Password, user.Password)

		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Incorrect email or password",
				"error":   true,
				"data":    nil,
			})
			return
		}

		if !user.IsEmailVerified {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			err = taskDistributor.DistributeTaskSendVerifyEmail(c, &worker.PayloadSendVerifyEmail{Email: body.Email}, opts...)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to distribute send verify email task",
					"error":   true,
					"data":    err,
				})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Email is not verified",
				"error":   true,
				"data":    nil,
			})
			return
		}

		// create token
		tokenMaker, err := auth.NewJwtMaker()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		accessToken, accessTokenPayload, err := tokenMaker.CreateToken(user.Id, config.AppConfig.ACCESS_TOKEN_DURATION)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		refreshToken, refreshTokenPayload, err := tokenMaker.CreateToken(user.Id, config.AppConfig.REFRESH_TOKEN_DURATION)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		session, err := dbService.CreateSession(c, db.CreateSessionRequest{
			Id:           refreshTokenPayload.ID,
			UserId:       user.Id,
			RefreshToken: refreshToken,
			ClientIp:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			ExpiresAt:    refreshTokenPayload.ExpiredAt,
		})

		if err != nil {
			me, ok := err.(*mysql.MySQLError)
			if ok && me.Number == MYSQL_KEY_EXISTS {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "session exists",
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

		user.Password = ""

		c.JSON(http.StatusOK, gin.H{
			"success":                  true,
			"message":                  "login successful",
			"error":                    nil,
			"user":                     user,
			"access_token":             accessToken,
			"refresh_token":            refreshToken,
			"access_token_expires_at":  accessTokenPayload.ExpiredAt,
			"refresh_token_expires_at": refreshTokenPayload.ExpiredAt,
			"session_id":               session.Id,
		})
	}
}
