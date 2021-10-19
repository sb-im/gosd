package api

import (
	"sb.im/gosd/luavm"
	"sb.im/gosd/state"
	"sb.im/gosd/storage"

	"github.com/go-oauth2/oauth2/v4/server"
)

type handler struct {
	cache   *state.State
	oauth   *server.Server
	store   *storage.Storage
	worker  *luavm.Worker
	baseURL string

	expireToken int
}
