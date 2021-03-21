package api

import (
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

func (h *handler) planRunning(w http.ResponseWriter, r *http.Request) {
	PlanID := request.RouteStringParam(r, "planID")

	service, ok := h.worker.Running[PlanID]
	if ok {
		json.OK(w, r, service.Task)
	} else {
		json.OK(w, r, nil)
	}
}

func (h *handler) planRunningDestroy(w http.ResponseWriter, r *http.Request) {
	PlanID := request.RouteStringParam(r, "planID")
	h.worker.Kill(PlanID)
	json.OK(w, r, []string{})
}
