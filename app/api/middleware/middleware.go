package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func RequestId() gin.HandlerFunc {
	return requestid.New(requestid.WithCustomHeaderStrKey("X-Trace-Id"))
}

func LogContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("traceid", requestid.Get(c))
		c.Next()
	}
}
