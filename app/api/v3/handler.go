package v3

import (
	"sb.im/gosd/app/config"
	"sb.im/gosd/app/service"

	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	srv *service.Service
	cfg *config.Config
}

func NewHandler(cfg *config.Config, orm *gorm.DB, srv *service.Service) *Handler {
	return &Handler{
		cfg: cfg,
		orm: orm,
		srv: srv,
	}
}
