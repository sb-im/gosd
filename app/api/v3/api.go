package v3

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"sb.im/gosd/app/service"
	"sb.im/gosd/luavm"
	"sb.im/gosd/app/docs"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

func NewApi(orm *gorm.DB, worker *luavm.Worker) http.Handler {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/gosd/api/v3"
	sr := r.Group("/gosd/api/v3")
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := NewHandler(orm, service.NewService(orm, worker))
	sr.GET("schedules", handler.scheduleIndex)
	sr.POST("schedules", handler.scheduleCreate)
	sr.PATCH("schedules/:id", handler.scheduleUpdate)
	sr.DELETE("schedules/:id", handler.scheduleDestroy)

	sr.GET("tasks", handler.TaskIndex)
	sr.POST("tasks", handler.TaskCreate)

	sr.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(200, gin.H{
			"message": "NO Router",
		})
	})
	return r
}
