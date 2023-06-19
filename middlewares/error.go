package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/httperr"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case httperr.ErrCustomError:
				c.AbortWithStatusJSON(e.StatusCode, map[string]string{"message": e.Message})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
