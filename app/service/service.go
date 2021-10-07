package service

import (
	"gorm.io/gorm"
	"sb.im/gosd/luavm"
)

type Service struct {
	orm    *gorm.DB
	worker *luavm.Worker
}

func NewService(orm *gorm.DB, worker *luavm.Worker) *Service {
	return &Service{orm, worker}
}
