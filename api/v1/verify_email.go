package v1

import (
	"net/http"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
)

func verifyEmail(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body db.VerifyEmailRequest
		err := c.ShouldBindQuery(&body)
		if err != nil {
			validationResponse(err, c)
			return
		}

		err = dbService.VerifyEmail(c, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "email verified",
			"error":   false,
		})
	}
}
