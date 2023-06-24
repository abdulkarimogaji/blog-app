package lambda

import (
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func uploadFileAPI(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "bad request",
				"message": err.Error(),
				"error":   true,
				"data":    nil,
			})
			return
		}

		response, err := cld.Upload.Upload(c, file, uploader.UploadParams{ResourceType: "auto"})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed to upload to cloudinary",
				"message": err.Error(),
				"error":   true,
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "upload successful",
			"error":   false,
			"data":    gin.H{"url": response.URL},
		})
	}
}
