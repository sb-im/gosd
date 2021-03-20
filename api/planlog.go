package api

import (
	"io/ioutil"
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
	PlanID := request.RouteStringParam(r, "planID")

	_, ok := h.worker.Running[PlanID]
	if ok {
		json.OK(w, r, []string{PlanID})
	} else {
		json.OK(w, r, []string{})
	}
}

func (h *handler) missionStop(w http.ResponseWriter, r *http.Request) {
	PlanID := request.RouteStringParam(r, "planID")
	h.worker.Kill(PlanID)
	json.OK(w, r, []string{})
}

func (h *handler) createPlanLog(w http.ResponseWriter, r *http.Request) {
	log := model.NewPlanLog()
	//log.LogID = request.RouteInt64Param(r, "logID")
	log.PlanID = request.RouteInt64Param(r, "planID")

	if mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err == nil && strings.HasPrefix(mediaType, "multipart/") {

		params, file, err := h.formData2Blob(r)
		if err != nil {
			json.ServerError(w, r, err)
			return
		}

		log.Attachments = file
		log.Extra = params
	}

	if err := h.store.CreatePlanLog(log); err != nil {
		json.ServerError(w, r, err)
		return
	}

	// Run Task
	if err := h.sendTask(log); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, log)
}

func (h *handler) sendTask(log *model.PlanLog) error {
	plan, err := h.store.PlanByID(log.PlanID)
	if err != nil {
		return err
	}

	var script []byte
	if blobID := plan.Files["lua"]; blobID != "" {
		int64ID, err := strconv.ParseInt(blobID, 10, 64)
		if err != nil {
			return err
		}

		blob, err := h.store.BlobByID(int64ID)
		if err != nil {
			return err
		}

		script, err = ioutil.ReadAll(blob.Reader)
		if err != nil {
			return err
		}
	}

	task := luavm.NewTask(log.ID, plan.NodeID, plan.ID)
	task.Files = plan.Files
	task.Extra = plan.Extra
	//task.URL   = h.baseURL + "/api/v1/plans/" + strconv.FormatInt(log.PlanID, 10) + "/logs/" + strconv.FormatInt(log.LogID, 10)
	task.URL = "api/v1/plans/%d?files=%s&token=%s"
	task.JobURL = "api/v1/plans/%d/jobs/%d?files=%s&token=%s"
	task.Script = script

	h.worker.Queue <- task

	return nil
}

func (h *handler) planLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.store.PlanLogs(request.RouteInt64Param(r, "planID"))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, logs)
}
