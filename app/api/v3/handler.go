package v3

import (
	"sb.im/gosd/app/service"

	"gorm.io/gorm"
)

type Handler struct {
	orm *gorm.DB
	srv *service.Service
}

func NewHandler(orm *gorm.DB, srv *service.Service) *Handler {
	return &Handler{
		orm: orm,
		srv: srv,
	}
}
