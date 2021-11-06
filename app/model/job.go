package model

type Job struct {
	Model

	TaskID int64  `json:"task_id"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
}
