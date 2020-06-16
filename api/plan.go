package api

import (
	"mime"
	"net/http"
	"strconv"
	"strings"

	"sb.im/gosd/model"

	"miniflux.app/http/response/json"
)

func (h *handler) createPlan(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	plan := model.NewPlan()
	if strings.HasPrefix(mediaType, "multipart/") {

		params, err := h.formData2Blob(r)
		if err != nil {
			json.ServerError(w, r, err)
			return
		}
		plan.Name = params["name"]
		plan.Description = params["description"]
		plan.NodeID, _ = strconv.ParseInt(params["node_id"], 10, 64)

		delete(params, "name")
		delete(params, "description")
		delete(params, "node_id")

		plan.Attachments = params

	} else {
		plan, err = decodePlanCreationPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
	}

	if err := h.store.CreatePlan(plan); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, plan)
}
