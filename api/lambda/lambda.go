package lambda

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(router *gin.RouterGroup, cloudinaryInstance *cloudinary.Cloudinary) {
	router.POST("/mail", sendMailAPI)
	router.POST("/upload", uploadFileAPI(cloudinaryInstance))
}
