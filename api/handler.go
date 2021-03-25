package api

import (
	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"
)

type handler struct {
	cache   *state.State
	store   *storage.Storage
	worker  *luavm.Worker
	baseURL string
}
