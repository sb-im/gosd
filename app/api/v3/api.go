package v3

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"sb.im/gosd/app/service"
	"sb.im/gosd/luavm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	handler := NewHandler(orm, service.NewService(orm, worker))

	// Init Auth Middleware
	handler.initAuth(sr)

	sr.GET("schedules", handler.scheduleIndex)
	sr.POST("schedules", handler.scheduleCreate)
	sr.PATCH("schedules/:id", handler.scheduleUpdate)
	sr.DELETE("schedules/:id", handler.scheduleDestroy)

	sr.POST("blobs", handler.blobCreate)
	sr.PUT("blobs", handler.blobUpdate)
	sr.PUT("blobs/:blobID", handler.blobUpdate)
	sr.GET("blobs/:blobID", handler.blobShow)

	sr.GET("tasks", handler.TaskIndex)
	sr.POST("tasks", handler.TaskCreate)

	sr.POST("teams", handler.TeamCreate)

	sr.POST("users", handler.UserCreate)
	sr.PATCH("users/:id", handler.UserUpdate)

	sr.GET("current", handler.current)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(200, gin.H{
			"message": "NO Router",
		})
	})

	handler.userOverride()
	return r
}
