package api

import (
	"log"
	"net/http"

	"github.com/abdulkarimogaji/blognado/api/lambda"
	v1 "github.com/abdulkarimogaji/blognado/api/v1"
	"github.com/abdulkarimogaji/blognado/config"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/abdulkarimogaji/blognado/worker"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

type Server interface {
	Start() error
}

type GinServer struct {
	DbService       db.DBService
	Router          *gin.Engine
	TaskDistributor worker.TaskDistributor
}

func NewServer(db db.DBService, taskDistributor worker.TaskDistributor) Server {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("none", func(fl validator.FieldLevel) bool { return true })
	}

	lambda.ConfigureRoutes(r.Group("/api/lambda/"))
	v1.ConfigureRoutes(r.Group("/v1/api/"), db, taskDistributor)

	r.GET("/health", func(c *gin.Context) {
		err := db.PingDB()

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
	return &GinServer{Router: r, DbService: db, TaskDistributor: taskDistributor}
}

func (s *GinServer) Start() error {
	return s.Router.Run(":" + config.AppConfig.PORT)
}
