package service

import (
	"sb.im/gosd/luavm"

	"gorm.io/gorm"
	"github.com/robfig/cron/v3"
)

type Service struct {
	orm    *gorm.DB
	worker *luavm.Worker
	// https://pkg.go.dev/github.com/robfig/cron/v3#Cron.Start
	cron   *cron.Cron
}

func NewService(orm *gorm.DB, worker *luavm.Worker) *Service {
	return &Service{orm, worker, cron.New()}
}
