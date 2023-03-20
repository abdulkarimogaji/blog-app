package v1

import (
	"database/sql"
	"net/http"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-sql-driver/mysql"
)

func createBlog(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body db.CreateBlogRequest
		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			validationResponse(err, c)
			return
		}

		blog, err := dbService.CreateBlog(body)
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

			me, ok := err.(*mysql.MySQLError)
			if ok && me.Number == MYSQL_KEY_EXISTS {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "blog slug already exists",
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
			"success": true,
			"message": "Blog created successfully",
			"data":    blog,
		})
	}
}
