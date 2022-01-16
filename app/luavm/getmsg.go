package luavm

import (
	"fmt"
)

const (
	topicNodeSys = "nodes/%s/%s"
	topicNodeMsg = "nodes/%s/msg/%s"
)

func (s Service) GetSys(id, msg string) (string, error) {
	return s.rdb.Get(s.ctx, fmt.Sprintf(topicNodeSys, id, msg)).Result()
}

func (s Service) GetMsg(id, msg string) (string, error) {
	return s.rdb.Get(s.ctx, fmt.Sprintf(topicNodeMsg, id, msg)).Result()
}
