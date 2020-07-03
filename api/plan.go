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
	user := h.helpCurrentUser(w, r)
	if user == nil {
		json.ServerError(w, r, errors.New("NotFound user"))
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	plan := model.NewPlan()
	if strings.HasPrefix(mediaType, "multipart/") {

		params, file, err := h.formData2Blob(r)
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

		plan.Attachments = file
		plan.Extra = params

	} else {
		plan, err = decodePlanCreationPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
	}

	// Add plan belong group
	plan.GroupID = user.ID

	if err := h.store.CreatePlan(plan); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, plan)
}

func (h *handler) plans(w http.ResponseWriter, r *http.Request) {
	user := h.helpCurrentUser(w, r)
	if user == nil {
		json.ServerError(w, r, errors.New("NotFound user"))
		return
	}

	plans, err := h.store.FindPlansByGroup(user.Group.ID)

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

func (h *handler) planUpdate(w http.ResponseWriter, r *http.Request) {
	planID := request.RouteInt64Param(r, "planID")
	plan, err := h.store.PlanByID(planID)
	if err != nil {
		json.BadRequest(w, r, errors.New("Unable to fetch this plan from the database"))
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}
	if strings.HasPrefix(mediaType, "multipart/") {

		params, file, err := h.formData2Blob(r)
		if err != nil {
			json.ServerError(w, r, err)
			return
		}
		plan.Name = params["name"]
		plan.Description = params["description"]
		plan.NodeID, _ = strconv.ParseInt(params["node_id"], 10, 64)

		// Only Update
		delete(params, "id")

		delete(params, "name")
		delete(params, "description")
		delete(params, "node_id")

		for key, value := range file {
			plan.Attachments[key] = value
		}

		for key, value := range params {
			plan.Extra[key] = value
		}

	} else {
		plan, err = decodePlanCreationPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
	}

	if err := h.store.UpdatePlan(plan); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, plan)
}

func (h *handler) planDestroy(w http.ResponseWriter, r *http.Request) {
	planID := request.RouteInt64Param(r, "planID")
	plan, err := h.store.PlanDestroy(planID)
	if err != nil {
		json.BadRequest(w, r, errors.New("Unable to fetch this plan from the database"))
		return
	}

	if plan == nil {
		json.NotFound(w, r)
		return
	}

	json.OK(w, r, plan)
}
