package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/requestid"
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
