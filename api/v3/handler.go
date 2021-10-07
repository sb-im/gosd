package v3

import (
	"gorm.io/gorm"
	"sb.im/gosd/storage"
)

type Handler struct {
	orm     *gorm.DB
	store   *storage.Storage
	baseURL string
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		orm: db,
	}
}
