package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		header := c.GetHeader("Authorization")
		if len(header) > 7 {
			tokenString = header[7:]
		}
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Access token not provided",
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}

		tokenMaker, err := NewJwtMaker()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"status":  "server error",
				"error":   true,
				"data":    nil,
			})
			return
		}
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": err.Error(),
				"status":  "unauthorized",
				"error":   true,
				"data":    nil,
			})
			return
		}
		c.Set("userId", payload.UserID)
		c.Next()
	}
}
