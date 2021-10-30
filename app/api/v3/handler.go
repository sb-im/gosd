package v3

import (
	"sb.im/gosd/app/service"

	"gorm.io/gorm"
)

var (
	DefaultConfig = &Config{
		StoragePath: "data/storage/",
	}
)

type Handler struct {
	orm *gorm.DB
	srv *service.Service
	cfg *Config
}

func NewHandler(orm *gorm.DB, srv *service.Service) *Handler {
	return &Handler{
		orm: orm,
		srv: srv,
		cfg: DefaultConfig,
	}
}

type Config struct {
	StoragePath string
}
