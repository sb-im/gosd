package luavm

import (
	"strconv"
)

type Task struct {
	id     int64             `json:"-"`
	nodeID int64             `json:"-"`
	planID int64             `json:"-"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
	URL    string            `json:"-"`
	Script []byte            `json:"-"`
	Job    Job               `json:"job"`
}

func NewTask(id, nodeID, planID int64) *Task {
	return &Task{
		id: id,
		nodeID: nodeID,
		planID: planID,
		Job: Job{
			JobID: id,
		},
	}
}

func (t *Task) NodeID() string {
	return strconv.FormatInt(t.nodeID, 10)
}

func (t *Task) PlanID() string {
	return strconv.FormatInt(t.planID, 10)
}

type Job struct {
	JobID  int64             `json:"job_id"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
}

func (t *Job) StringJobID() string {
	return strconv.FormatInt(t.JobID, 10)
}

