package api

import (
	"log"
	"net/http"

	"github.com/abdulkarimogaji/blognado/api/lambda"
	v1 "github.com/abdulkarimogaji/blognado/api/v1"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

func RunServer() error {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("none", func(fl validator.FieldLevel) bool { return true })
	}

	lambda.ConfigureRoutes(r.Group("/api/lambda/"))
	v1.ConfigureRoutes(r.Group("/v1/api/"))

	r.GET("/health", func(c *gin.Context) {
		err := db.DbConn.Ping()

		if err != nil {
			log.Println("error here", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to ping the database",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r.Run(":" + config.AppConfig.PORT)
}
