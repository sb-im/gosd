package api

import (
	"errors"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"sb.im/gosd/model"

	"miniflux.app/http/request"
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

func (h *handler) plans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.store.Plans()

	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.OK(w, r, plans)
}

func (h *handler) planByID(w http.ResponseWriter, r *http.Request) {
	planID := request.RouteInt64Param(r, "planID")
	plan, err := h.store.PlanByID(planID)
	if err != nil {
		json.BadRequest(w, r, errors.New("Unable to fetch this plan from the database"))
		return
	}

	if plan == nil {
		json.NotFound(w, r)
		return
	}

	//user.UseTimezone(request.UserTimezone(r))
	json.OK(w, r, plan)
}
