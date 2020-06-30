package api

import (
	"io/ioutil"
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

func (h *handler) actionHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("data/" + request.RouteStringParam(r, "action") + ".json")
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}
