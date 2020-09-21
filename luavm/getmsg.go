package luavm

import (
	"errors"
)

func (s *Service) GetMsg(id, msg string) (string, error) {
	raw, err := s.State.NodeGet(id, msg)
	return string(raw), err
}

func (s *Service) GetStatus(id string) (string, error) {
	if data := s.State.Node[id]; data != nil {
		return string(data.Status.Raw), nil
	}
	return "", errors.New("No Find Status")
}
