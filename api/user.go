package api

import (
	"errors"
	"net/http"
	"strings"

	"sb.im/gosd/auth"
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
	// Single user mode
	// User.ID == 1
	if auth.AuthMethod == auth.NoAuth {
		user, _ := h.store.UserByID(1)
		return user
	}
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		return h.store.GetCurrentUser(token[1])
	}
	return nil
}
