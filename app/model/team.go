package model

type Team struct {
	Model

	Name string `json:"name" form:"name"`
}
