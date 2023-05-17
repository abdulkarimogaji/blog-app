package v1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func refreshToken(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body refreshTokenRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)

		// request validation
		if err != nil {
			validationResponse(err, c)
			return
		}

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

		payload, err := tokenMaker.VerifyToken(body.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		session, err := dbService.GetSession(c, payload.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "session not found",
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

		if session.IsBlocked {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "session blocked",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		if session.UserId != payload.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Incorrect session user",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		if session.RefreshToken != body.RefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "mismatched session token",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		if time.Now().After(session.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "session expired",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		// create new token
		accessToken, accessTokenPayload, err := tokenMaker.CreateToken(session.UserId, config.AppConfig.ACCESS_TOKEN_DURATION)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":            accessToken,
			"access_token_expires_at": accessTokenPayload.ExpiredAt,
			"session_id":              session.Id,
		})
	}
}
