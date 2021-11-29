package middlewares

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		requiredToken := c.Request.Header["Authorization"]

		if len(requiredToken) == 0 {
			responseWithError(c, 403, "Please login to your account")
		}

		token := strings.Replace(requiredToken[0], "Bearer ", "", 1)
		if token != os.Getenv("APP_TOKEN") {
			responseWithError(c, 403, "Invalid token")
		}
	}
}
