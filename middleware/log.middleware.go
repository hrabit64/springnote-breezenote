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

		fmt.Printf("[%d] %s %s Content-Type: %s  In %d ns IP: %s\n",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Header().Get("Content-Type"),
			latency.Nanoseconds(),
			c.ClientIP(),
		)
	}
}
