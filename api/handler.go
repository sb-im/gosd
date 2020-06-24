package api

import (
	"sb.im/gosd/luavm"
	"sb.im/gosd/storage"
)

type handler struct {
	store   *storage.Storage
	worker  *luavm.Worker
	baseURL string
}
