package model

type Group struct {
	ID    int64             `json:"id"`
	Name  string            `json:"name"`
	Extra map[string]string `json:"extra"`
}

func NewGroup() *Group {
	return &Group{}
}
