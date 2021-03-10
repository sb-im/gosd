package state

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	redis "github.com/gomodule/redigo/redis"
)

type State struct {
	Plan map[int]PlanState
	Mqtt mqtt.Client
	Conn redis.Conn
}

type PlanState struct {
	id int64
}

func NewState(rawURL string) *State {
	//c, err := redis.DialURL(os.Getenv("REDIS_URL"))
	c, err := redis.DialURL(rawURL)
	if err != nil {
		// TODO: handle connection error
	}
	//defer c.Close()

	return &State{
		Conn: c,
	}
}

func (s *State) Record(key string, value []byte) error {
	if _, err := s.Conn.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

func (s *State) GetNode(id, msg string) (string, error) {
	return redis.String(s.Conn.Do("GET", fmt.Sprintf("nodes/%s/%s", id, msg)))
}

func (s *State) GetNodeMsg(id, msg string) ([]byte, error) {
	return redis.Bytes(s.Conn.Do("GET", fmt.Sprintf("nodes/%s/msg/%s", id, msg)))
}
