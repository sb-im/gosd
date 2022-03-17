package v3

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	rdb *redis.Client
	srv *service.Service
	cfg *config.Config
	ofs *storage.Storage
}

func NewHandler(cfg *config.Config, orm *gorm.DB, rdb *redis.Client, srv *service.Service, ofs *storage.Storage) *Handler {
	return &Handler{
		cfg: cfg,
		orm: orm,
		rdb: rdb,
		srv: srv,
		ofs: ofs,
	}
}
