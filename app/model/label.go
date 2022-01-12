package model

type Label struct {
	Model

	Name   string `json:"name"`
	TeamID uint   `json:"-"`
	Files  JSON   `json:"files"`
	Extra  JSON   `json:"extra"`
}
