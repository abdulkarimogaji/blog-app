package v1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func login(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body db.LoginRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)

		// request validation
		if err != nil {
			validationResponse(err, c)
			return
		}

		// get user
		user, err := dbService.GetUserByEmail(body.Email)
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

		// create token
		maker, err := auth.NewJwtMaker()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		token, err := maker.CreateToken(user.Id, time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
				"data":    nil,
			})
			return
		}

		user.Password = ""

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "login successful",
			"error":   nil,
			"user":    user,
			"token":   token,
		})
	}
}
