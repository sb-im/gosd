package v3

import (
	"sb.im/gosd/app/service"

	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	srv *service.Service
	cfg *Config
}

func NewHandler(cfg *Config, orm *gorm.DB, srv *service.Service) *Handler {
	return &Handler{
		cfg: cfg,
		orm: orm,
		srv: srv,
	}
}
