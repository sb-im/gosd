package service

import (
	"sb.im/gosd/app/luavm"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Service struct {
	orm    *gorm.DB
	rdb    *redis.Client
	worker *luavm.Worker
	// https://pkg.go.dev/github.com/robfig/cron/v3#Cron.Start
	cron *cron.Cron
}

func NewService(orm *gorm.DB, rdb *redis.Client, worker *luavm.Worker) *Service {
	return &Service{orm, rdb, worker, cron.New()}
}
