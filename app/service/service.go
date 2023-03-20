package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Service struct {
	orm *gorm.DB
	rdb *redis.Client
	// https://pkg.go.dev/github.com/robfig/cron/v3#Cron.Start
	cron *cron.Cron
}

func NewService(orm *gorm.DB, rdb *redis.Client) *Service {
	return &Service{orm, rdb, cron.New()}
}
