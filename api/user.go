package api

import (
	"errors"
	"net/http"
	"strings"

	"sb.im/gosd/model"

	"miniflux.app/http/response/json"
)

func (h *handler) currentUser(w http.ResponseWriter, r *http.Request) {
	user := h.helpCurrentUser(w, r)
	if user == nil {
		json.ServerError(w, r, errors.New("No login"))
		return
	}

	json.OK(w, r, user)
}

func (h *handler) helpCurrentUser(w http.ResponseWriter, r *http.Request) *model.User {
	return h.store.GetCurrentUser(strings.Split(r.Header.Get("Authorization"), " ")[1])
}
