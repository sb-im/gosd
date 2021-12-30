package v3

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"sb.im/gosd/app/api/v3"
	"sb.im/gosd/app/service"
	"sb.im/gosd/luavm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NewApi(orm *gorm.DB, worker *luavm.Worker) http.Handler {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	sr := r.Group("/gosd/api/v3")
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := v3.NewHandler(orm, service.NewService(orm, redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}), worker))

	// Init Auth Middleware
	handler.InitAuth(sr)

	sr.GET("schedules", handler.ScheduleIndex)
	sr.POST("schedules", handler.ScheduleCreate)
	sr.PATCH("schedules/:id", handler.ScheduleUpdate)
	sr.DELETE("schedules/:id", handler.ScheduleDestroy)
	sr.POST("schedules/:id/toggle", handler.ScheduleToggle)

	sr.POST("blobs", handler.BlobCreate)
	sr.PUT("blobs", handler.BlobUpdate)
	sr.PUT("blobs/:blobID", handler.BlobUpdate)
	sr.GET("blobs/:blobID", handler.BlobShow)

	sr.POST("tasks", handler.TaskCreate)
	sr.GET("tasks", handler.TaskIndex)
	sr.GET("tasks/:id", handler.TaskShow)
	sr.PUT("tasks/:id", handler.TaskUpdate)
	sr.DELETE("tasks/:id", handler.TaskDestroy)

	sr.POST("teams", handler.TeamCreate)
	sr.POST("teams/users", handler.TeamUserAdd)

	sr.POST("users", handler.UserCreate)
	sr.PATCH("users/:id", handler.UserUpdate)

	sr.POST("test", handler.MqttUserCreate)

	sr.GET("current", handler.Current)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(200, gin.H{
			"message": "NO Router",
		})
	})

	handler.InitSeed()
	handler.UserOverride()
	return r
}
