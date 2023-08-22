package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-starter/pkg/jwt"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.TokenValid(c); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, map[string]any{})
			c.Abort()
			return
		}
		c.Next()
	}
}
