package v3

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	rdb *redis.Client
	srv *service.Service
	cfg *config.Config
	ofs *storage.Storage
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		cfg: srv.Cfg(),
		orm: srv.Orm(),
		rdb: srv.Rdb(),
		ofs: srv.Ofs(),
		srv: srv,
	}
}
