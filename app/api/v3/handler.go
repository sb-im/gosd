package v3

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"
	"sb.im/gosd/app/storage"

	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	srv *service.Service
	cfg *config.Config
	ofs *storage.Storage
}

func NewHandler(cfg *config.Config, orm *gorm.DB, srv *service.Service, ofs *storage.Storage) *Handler {
	return &Handler{
		cfg: cfg,
		orm: orm,
		srv: srv,
		ofs: ofs,
	}
}
