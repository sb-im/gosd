package model

type Plan2 struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	NodeID      int64             `json:"node_id"`
	GroupID     int64             `json:"-"`
	Files       map[string]string `json:"files"`
	Extra       map[string]string `json:"extra"`
	RecordTime
}

func NewPlan2() *Plan2 {
	return &Plan2{
		Files: make(map[string]string),
		Extra: make(map[string]string),
	}
}

type Plan2s []*Plan2
