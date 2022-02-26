package api

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"sb.im/gosd/app/api/v3"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewApi(cfg *config.Config, orm *gorm.DB, srv *service.Service) http.Handler {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	sr := r.Group("/gosd/api/v3")
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := v3.NewHandler(cfg, orm, srv)

	// Init Auth Middleware
	handler.InitAuth(sr)

	sr.GET("/status", handler.Status)

	// === Manager API { ===
	sr.GET("/teams", handler.TeamIndex)
	sr.POST("/teams", handler.TeamCreate)
	sr.PATCH("/teams/:id", handler.TeamUpdate)

	// This api disable
	// sr.POST("/teams/users", handler.TeamUserAdd)

	// === Manager API } ===

	sr.GET("schedules", handler.ScheduleIndex)
	sr.POST("schedules", handler.ScheduleCreate)
	sr.PATCH("schedules/:id", handler.ScheduleUpdate)
	sr.DELETE("schedules/:id", handler.ScheduleDestroy)
	sr.POST("/schedules/:id/trigger", handler.ScheduleTrigger)

	sr.POST("blobs", handler.BlobCreate)
	sr.PUT("blobs", handler.BlobUpdate)
	sr.PUT("blobs/:blobID", handler.BlobUpdate)
	sr.GET("blobs/:blobID", handler.BlobShow)

	sr.POST("nodes", handler.NodeCreate)
	sr.GET("nodes", handler.NodeIndex)
	sr.GET("nodes/:id", handler.NodeShow)
	sr.PUT("nodes/:id", handler.NodeUpdate)
	sr.DELETE("nodes/:id", handler.NodeDestroy)

	sr.POST("tasks", handler.TaskCreate)
	sr.GET("tasks", handler.TaskIndex)
	sr.GET("tasks/:id", handler.TaskShow)
	sr.PUT("tasks/:id", handler.TaskUpdate)
	sr.DELETE("tasks/:id", handler.TaskDestroy)

	sr.GET("tasks/:id/jobs", handler.JobIndex)
	sr.POST("tasks/:id/jobs", handler.JobCreate)

	sr.GET("/users", handler.UserIndex)
	sr.POST("users", handler.UserCreate)
	sr.POST("/users/:user_id/teams/:team_id", handler.UserAddTeam)
	sr.PATCH("users/:id", handler.UserUpdate)

	sr.POST("mqtt/url", handler.MqttUserCreate)

	sr.GET("current", handler.Current)

	r.NoRoute(func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "NO Router",
		})
	})

	handler.InitSeed()
	handler.UserOverride()
	return r
}
