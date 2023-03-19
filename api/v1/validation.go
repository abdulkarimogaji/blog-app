package v1

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "json":
		return "Should be a valid json"
	case "min":
		return toSnakeCase(fe.Field()) + " should be at least " + fe.Param() + " characters"
	case "url":
		return "should be a valid url"
	case "datetime":
		return "should be in format: " + fe.Param()
	case "email":
		return "should be a valid email address"
	}
	return "Unknown error"
}

func validationResponse(err error, c *gin.Context) {

	if err == io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "empty body"})
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]errorMsg, len(ve))
		for i, fe := range ve {
			out[i] = errorMsg{toSnakeCase(fe.Field()), getErrorMsg(fe)}
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "validation error", "errors": out})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
