package service

import (
	"testing"
)

func TestCall(t *testing.T) {
	srv := NewJsonService(NewService(nil, nil, nil))
	srv.Call("taskRun", []byte(`{"id": 1}`))
}
