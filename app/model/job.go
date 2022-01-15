package model

type Job struct {
	Model

	TaskID uint `json:"-"`
	Task   Task `json:"-"`
	Files  JSON `json:"files"`
	Extra  JSON `json:"extra"`
}
