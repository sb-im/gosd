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

func (h *handler) createPlan2(w http.ResponseWriter, r *http.Request) {
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
	plan2 := model.NewPlan2()
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
		plan2, err = decodePlan2CreationPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
		copyPlan2Plan(plan2, plan)
	}

	// Add plan belong group
	plan.GroupID = user.ID

	if err := h.store.CreatePlan(plan); err != nil {
		json.ServerError(w, r, err)
		return
	}

	plan2.ID = plan.ID

	json.Created(w, r, plan2)
}

func copyPlan2Plan(plan2 *model.Plan2, plan *model.Plan) {
	plan.Name = plan2.Name
	plan.Description = plan2.Description
	plan.NodeID = plan2.NodeID
	plan.Attachments = plan2.Files
	plan.Extra = plan2.Extra
}

func (h *handler) plan2s(w http.ResponseWriter, r *http.Request) {
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

	plan2s := make(model.Plan2s, len(plans))
	for i, p := range plans {
		plan2s[i] = model.NewPlan2()
		plan2s[i].ID = p.ID
		plan2s[i].Name = p.Name
		plan2s[i].Description = p.Description
		plan2s[i].NodeID = p.NodeID
		plan2s[i].Files = p.Attachments
		plan2s[i].Extra = p.Extra
	}

	json.OK(w, r, plan2s)
}

func (h *handler) updatePlan2(w http.ResponseWriter, r *http.Request) {
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

	plan2 := model.NewPlan2()
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
		plan2, err = decodePlan2CreationPayload(r.Body)
		if err != nil {
			json.BadRequest(w, r, err)
			return
		}
		copyPlan2Plan(plan2, plan)
	}

	if err := h.store.UpdatePlan(plan); err != nil {
		json.ServerError(w, r, err)
		return
	}

	plan2.ID = plan.ID

	json.OK(w, r, plan2)
}
