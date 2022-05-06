package v3

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"
	"sb.im/gosd/app/store"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	rdb *redis.Client
	srv *service.Service
	cfg *config.Config
	ofs *storage.Storage

	store *store.Store
}

func NewHandler(s *store.Store, srv *service.Service) *Handler {
	return &Handler{
		cfg: s.Cfg(),
		orm: s.Orm(),
		rdb: s.Rdb(),
		ofs: s.Ofs(),
		srv: srv,

		store: s,
	}
}
