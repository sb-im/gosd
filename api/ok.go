package api

import (
	"net/http"
	"time"

	"miniflux.app/http/response/json"
)

func (h *handler) ok(w http.ResponseWriter, r *http.Request) {
	json.OK(w, r, struct {
		Ok string `json:"ok"`
	}{
		Ok: time.Now().String(),
	})
}
