package api

import (
	"errors"
	"net/http"
	"strings"

	"miniflux.app/http/response/json"
)

func (h *handler) currentUser(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.Header.Get("Authorization"), " ")[1]
	user := h.store.GetCurrentUser(key)

	if user == nil {
		json.ServerError(w, r, errors.New("No login"))
		return
	}

	json.OK(w, r, user)
}
