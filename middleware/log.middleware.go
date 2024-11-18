package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		t := time.Now()

		c.Next()

		latency := time.Since(t)

		fmt.Printf("[%d] %s Content-Type: %s Url :%s In %d\n",
			c.Writer.Status(),
			c.Request.Method,
			c.Writer.Header().Get("Content-Type"),
			c.Request.RequestURI,
			latency,
		)
	}
}
