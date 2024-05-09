package middleware

import (
	"log"
	"molta9a/server/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorToText(e string) string {
	switch e {
	case "required":
		return utils.InvalidFieldRequired.Error()
	case "max":
		return utils.InvalidMaxLength.Error()
	case "min":
		return utils.InvalidMinLength.Error()
	case "email":
		return utils.InvalidEmail.Error()
	case "len":
		return utils.InvalidLength.Error()
	}
	return utils.InvalidField.Error()
}

// This method collects all errors
func Errors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"error": e.Error()})
					}
				case gin.ErrorTypeBind:
					errs := e.Err.(validator.ValidationErrors)
					list := make(map[string]string)
					for _, err := range errs {
						list[strings.ToLower(err.Field())] = ValidationErrorToText(err.Tag())
					}

					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, gin.H{"errors": list})

				default:
					// Log all errors
					log.Println(e.Err)
				}

			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
			}
		}
	}
}
