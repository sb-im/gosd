package service

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/storage"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Service struct {
	cfg *config.Config
	orm *gorm.DB
	rdb *redis.Client
	ofs *storage.Storage
	// https://pkg.go.dev/github.com/robfig/cron/v3#Cron.Start
	cron *cron.Cron
}

func NewService(cfg *config.Config, orm *gorm.DB, rdb *redis.Client, ofs *storage.Storage) *Service {
	return &Service{cfg, orm, rdb, ofs, cron.New()}
}

func (s *Service) Cfg() *config.Config {
	return s.cfg
}

func (s *Service) Orm() *gorm.DB {
	return s.orm
}

func (s *Service) Rdb() *redis.Client {
	return s.rdb
}

func (s *Service) Ofs() *storage.Storage {
	return s.ofs
}
