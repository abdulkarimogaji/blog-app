package lambda

import "github.com/gin-gonic/gin"

func ConfigureRoutes(router *gin.RouterGroup) {
	router.POST("/mail", sendMailAPI)
}
