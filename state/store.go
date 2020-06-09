package state

import (
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type State struct {
	Plan map[int]PlanState
	Node map[string]*NodeState
	Mqtt mqtt.Client
}

type PlanState struct {
	id int64
}

type NodeState struct {
	Msg map[string][]byte
}

func NewState() *State {
	return &State{
		Node: make(map[string]*NodeState),
	}
}

func (s *State) NodePut(id, msg string, payload []byte) error {
	if s.Node[id] == nil {
		s.Node[id] = &NodeState{
			Msg: map[string][]byte{},
		}
	}
	s.Node[id].Msg[msg] = payload

	return nil
}

func (s *State) NodeGet(id, msg string) ([]byte, error) {
	if s.Node[id] == nil {
		return []byte{}, errors.New("No Find this Message For Id: " + id)
	}

	if s.Node[id].Msg[msg] == nil {
		return []byte{}, errors.New("No Find this Message: " + msg)
	}

	return s.Node[id].Msg[msg], nil
}
