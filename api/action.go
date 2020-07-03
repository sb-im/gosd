package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

func (h *handler) actionHandler(w http.ResponseWriter, r *http.Request) {
	user := h.helpCurrentUser(w, r)
	if user != nil {
		content, err := ioutil.ReadFile("data/" + strconv.FormatInt(user.Group.ID, 10) + "/" + request.RouteStringParam(r, "action") + ".json")

		if err != nil {
			fmt.Println(err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
			return
		}

	}

	content, err := ioutil.ReadFile("data/" + request.RouteStringParam(r, "action") + ".json")
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}
