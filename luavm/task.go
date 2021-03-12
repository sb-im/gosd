package luavm

type Task struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	NodeID string            `json:"node_id"`
	PlanID string            `json:"-"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
	URL    string            `json:"-"`
	Script []byte            `json:"-"`
}
