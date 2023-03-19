package v1

import (
	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/abdulkarimogaji/blognado/db"
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(router *gin.RouterGroup, db db.DBService) {
	router.POST("/signup", signUp(db))
	router.POST("/login", login(db))
	router.Group("/", auth.AuthorizeClient())
}
