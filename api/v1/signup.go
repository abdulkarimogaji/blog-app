package v1

import (
	"net/http"
	"time"

	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-playground/validator/v10"
)

type signUpRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func signUp(c *gin.Context) {
	var body signUpRequest
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "bad request",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	hashedPassword, err := util.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}
	stmt, err := db.DbConn.Prepare("INSERT INTO user (first_name, last_name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	createdAt := time.Now()
	updatedAt := time.Now()

	result, err := stmt.Exec(body.FirstName, body.LastName, body.Email, hashedPassword, createdAt, updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": int(id),
		"error":   false,
	})
}
