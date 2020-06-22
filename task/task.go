package task

type Task interface {
	ID() string
	LogID() string
	NodeID() string
	URL() string
	Script() []byte
}
