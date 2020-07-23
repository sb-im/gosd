package luavm

type Task struct {
	PlanID string
	Attach map[string]string
	Extra  map[string]string
	NodeID string
	URL    string
	Script []byte
}
