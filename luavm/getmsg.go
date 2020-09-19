package luavm

import (
	"errors"
	"strconv"
)

func (s *Service) GetMsg(id, msg string) (string, error) {
	raw, err := s.State.NodeGet(id, msg)
	return string(raw), err
}

func (s *Service) GetID(id, str string) (string, error) {
	if n := s.State.Node[id]; n != nil && str == "link_id" {
		return strconv.Itoa(n.Status.GetID("")), nil
	}
	return "", errors.New("No Find id")
}

func (s *Service) GetStatus(id string) (string, error) {
	if data := s.State.Node[id]; data != nil {
		return string(data.Status.Raw), nil
	}
	return "", errors.New("No Find Status")
}
