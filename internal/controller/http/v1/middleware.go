package v1

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// UUIDParamValidator checks if the specified path parameters match the UUID format.
func UUIDParamValidator(paramNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, name := range paramNames {
			val := c.Param(name)
			if val != "" && !uuidRegex.MatchString(val) {
				errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid UUID format for parameter: "+name, nil)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
