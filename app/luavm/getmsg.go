package luavm

import (
	"context"
	"fmt"
)

const (
	topicNodeSys = "nodes/%s/%s"
	topicNodeMsg = "nodes/%s/msg/%s"
)

func (s Service) GetSys(id, msg string) (string, error) {
	return s.rdb.Get(context.Background(), fmt.Sprintf(topicNodeSys, id, msg)).Result()
}

func (s Service) GetMsg(id, msg string) (string, error) {
	return s.rdb.Get(context.Background(), fmt.Sprintf(topicNodeMsg, id, msg)).Result()
}

// TODO: Deprecated
func (s Service) GetStatus(id string) (string, error) {
	return s.GetSys(id, "status")
}

// TODO: Deprecated
func (s Service) GetNetwork(id string) (string, error) {
	return s.GetSys(id, "network")
}
