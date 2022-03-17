package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"sb.im/gosd/app/api/v3"
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var ApiPrefix = "/gosd/api/v3"

func NewApi(cfg *config.Config, orm *gorm.DB, rdb *redis.Client, srv *service.Service, ofs *storage.Storage) http.Handler {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	if u, err := url.Parse(cfg.BaseURL); err == nil {
		ApiPrefix = u.Path
	}

	sr := r.Group(ApiPrefix)
	sr.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := v3.NewHandler(cfg, orm, rdb, srv, ofs)

	// Init Auth Middleware
	if err := v3.InitAuthMiddleware(sr, handler); err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	sr.GET("/status", handler.Status)

	// === Manager API { ===
	sr.GET("/teams", handler.TeamIndex)
	sr.POST("/teams", handler.TeamCreate)
	sr.PATCH("/teams/:id", handler.TeamUpdate)

	// This api disable
	// sr.POST("/teams/users", handler.TeamUserAdd)

	sr.GET("/users", handler.UserIndex)
	sr.POST("/users", handler.UserCreate)
	sr.PATCH("/users/:id", handler.UserUpdate)

	sr.POST("/users/:user_id/teams/:team_id", handler.UserAddTeam)
	// === Manager API } ===

	sr.GET("schedules", handler.ScheduleIndex)
	sr.POST("schedules", handler.ScheduleCreate)
	sr.PATCH("schedules/:id", handler.ScheduleUpdate)
	sr.DELETE("schedules/:id", handler.ScheduleDestroy)
	sr.POST("/schedules/:id/trigger", handler.ScheduleTrigger)

	sr.POST("/blobs", handler.BlobCreate)
	sr.GET("/blobs/:blobID", handler.BlobShow)

	// blob cannot be updated
	// sr.PUT("/blobs", handler.BlobUpdate)
	// sr.PUT("/blobs/:blobID", handler.BlobUpdate)

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

	sr.POST("/tasks/:id/running", handler.TaskRunningCreate)
	sr.DELETE("/tasks/:id/running", handler.TaskRunningDestroy)

	sr.GET("tasks/:id/jobs", handler.JobIndex)
	sr.POST("tasks/:id/jobs", handler.JobCreate)

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
