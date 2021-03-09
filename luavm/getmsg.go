package luavm

import (
	"errors"
)

// TODO: Need Improve
func (s *Service) GetMsg(id, msg string) (string, error) {
	raw, err := s.State.NodeGet(id, msg)
	return string(raw), err
}

// TODO: Need Fix this
func (s *Service) GetStatus(id string) (string, error) {
	if data := s.State.Node[id]; data != nil {
		return string(data.Status.Raw), nil
	}
	return "", errors.New("No Find Status")
}
