package luavm

import (
	"strconv"
)

type Task struct {
	ID     int64             `json:"-"`
	NodeID int64             `json:"-"`
	PlanID int64             `json:"-"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
	URL    string            `json:"-"`
	JobURL string            `json:"-"`
	Script []byte            `json:"-"`
	Job    Job               `json:"job"`
}

func NewTask(id, nodeID, planID int64) *Task {
	return &Task{
		ID: id,
		NodeID: nodeID,
		PlanID: planID,
		Job: Job{
			JobID: id,
		},
	}
}

func (t *Task) StringNodeID() string {
	return strconv.FormatInt(t.NodeID, 10)
}

func (t *Task) StringPlanID() string {
	return strconv.FormatInt(t.PlanID, 10)
}

type Job struct {
	JobID  int64             `json:"job_id"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
}

func (t *Job) StringJobID() string {
	return strconv.FormatInt(t.JobID, 10)
}

