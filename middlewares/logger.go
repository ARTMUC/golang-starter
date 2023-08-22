package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.FullPath())
		fmt.Println(c.HandlerNames())
		c.Next()
	}
}
