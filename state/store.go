package state

type State struct {
	Plan map[int]PlanState
	Node map[int]NodeState
}

type PlanState struct {
	id int64
}

type NodeState struct {
	Msg map[string][]byte
}
