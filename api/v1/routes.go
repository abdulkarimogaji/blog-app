package v1

import (
	"github.com/abdulkarimogaji/blognado/api/middleware/auth"
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(router *gin.RouterGroup) {
	router.POST("/signup", signUp)
	router.POST("/login", login)
	router.Group("/", auth.AuthorizeClient())
}
