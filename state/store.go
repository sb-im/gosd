package state

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	redis "github.com/gomodule/redigo/redis"
)

type State struct {
	Plan map[int]PlanState
	Node map[string]*NodeState
	Mqtt mqtt.Client
	Conn redis.Conn
}

type PlanState struct {
	id int64
}

type NodeState struct {
	Status NodeStatus
	Msg    map[string][]byte
}

func NewState() *State {
	//c, err := redis.DialURL(os.Getenv("REDIS_URL"))
	c, err := redis.DialURL("redis://localhost:6379/0")
	if err != nil {
		// TODO: handle connection error
	}
	//defer c.Close()

	return &State{
		Node: make(map[string]*NodeState),
		Conn: c,
	}
}

func (s *State) Record(key string, value []byte) error {
	if _, err := s.Conn.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

func (s *State) NodeGet(id, msg string) ([]byte, error) {
	return redis.Bytes(s.Conn.Do("GET", fmt.Sprintf("nodes/%s/msg/%s", id, msg)))
}

func (s *State) addNode(id string) error {
	s.Node[id] = &NodeState{
		Msg: map[string][]byte{},
	}
	return nil
}

func (s *State) SetNodeStatus(id string, payload []byte) error {
	if s.Node[id] == nil {
		s.addNode(id)
	}
	return s.Node[id].Status.SetStatus(payload)
}
