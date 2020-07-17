package luavm

type Task struct {
	PlanID string
	NodeID string
	URL    string
	Script []byte
}
