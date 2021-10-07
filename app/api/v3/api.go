package v3

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"sb.im/gosd/app/service"
	"sb.im/gosd/luavm"

	"github.com/gin-gonic/gin"
)

func NewApi(orm *gorm.DB, worker *luavm.Worker) http.Handler {
	r := gin.Default()
	sr := r.Group("/gosd/api/v3")
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := NewHandler(orm, service.NewService(orm, worker))
	sr.GET("schedule", handler.ScheduleIndex)
	sr.POST("schedule", handler.ScheduleCreate)
	sr.POST("schedule/:id/toggle", handler.ScheduleToggle)

	sr.GET("tasks", handler.TaskIndex)
	sr.POST("tasks", handler.TaskCreate)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(200, gin.H{
			"message": "NO Router",
		})
	})
	return r
}
