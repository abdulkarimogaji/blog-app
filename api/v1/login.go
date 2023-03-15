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

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token"`
}

func login(c *gin.Context) {
	var body loginRequest
	err := c.ShouldBindBodyWith(&body, binding.JSON)

	// request validation
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "bad request",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	// get user
	var resp loginResponse
	row := db.DbConn.QueryRow("SELECT id, first_name, last_name, password, email, created_at, updated_at from user WHERE email = ?", body.Email)
	err = row.Scan(&resp.Id, &resp.FirstName, &resp.LastName, &resp.Password, &resp.Email, &resp.CreatedAt, &resp.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "not found",
				"message": "User not found",
				"error":   true,
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	//check password
	ok := util.VerifyPassword(body.Password, resp.Password)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "bad request",
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
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	token, err := maker.CreateToken(resp.Id, time.Minute*5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "server error",
			"message": err.Error(),
			"error":   true,
			"data":    nil,
		})
		return
	}

	resp.Token = token
	resp.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
		"error":   false,
		"data":    resp,
	})
}
