package api

import (
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"sb.im/gosd/luavm"
	"sb.im/gosd/model"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
)

func (h *handler) missionQueue(w http.ResponseWriter, r *http.Request) {
	json.OK(w, r, []string{})
}

func (h *handler) createPlanLog(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	log := model.NewPlanLog()
	//log.LogID = request.RouteInt64Param(r, "logID")
	log.PlanID = request.RouteInt64Param(r, "planID")

	if strings.HasPrefix(mediaType, "multipart/") {

		params, file, err := h.formData2Blob(r)
		if err != nil {
			json.ServerError(w, r, err)
			return
		}

		log.Attachments = file
		log.Extra = params

	} else {
		fmt.Println("Not Form Data")
	}

	if err := h.store.CreatePlanLog(log); err != nil {
		json.ServerError(w, r, err)
		return
	}

	plan, err := h.store.PlanByID(log.PlanID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	h.worker.Queue <- &luavm.Task{
		NodeID: strconv.FormatInt(plan.NodeID, 10),
		URL:    "1/12/3/4/4",
		Script: []byte{},
	}

	json.Created(w, r, log)
}
