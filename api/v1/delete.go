package v1

import (
	"net/http"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
)

type deleteParams struct {
	TableName string `uri:"tableName" binding:"required"`
	Id        int    `uri:"id" binding:"required,number"`
}

func deleteHandler(dbService db.DBService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params deleteParams
		err := c.ShouldBindUri(&params)
		if err != nil {
			validationResponse(err, c)
			return
		}

		id, err := dbService.DeleteRow(params.TableName, params.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "server error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": id,
		})
	}
}
