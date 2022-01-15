package model

type Job struct {
	Model

	TaskID uint `json:"task_id"`
	Task   Task `json:"team"`
	Files  JSON `json:"files"`
	Extra  JSON `json:"extra"`
}
