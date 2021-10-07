package api

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"sb.im/gosd/api/v3"

	"github.com/gin-gonic/gin"
)

func v3Handler(orm *gorm.DB) http.Handler {
	r := gin.Default()
	sr := r.Group("/gosd/api/v3")
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := v3.NewHandler(orm)
	sr.GET("schedule", handler.ScheduleIndex)
	sr.POST("schedule", handler.ScheduleCreate)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(200, gin.H{
			"message": "NO Router",
		})
	})
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	return r
}
