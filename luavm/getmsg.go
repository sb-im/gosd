package luavm

// TODO: Need Improve
func (s *Service) GetMsg(id, msg string) (string, error) {
	raw, err := s.State.NodeGet(id, msg)
	return string(raw), err
}

func (s *Service) GetStatus(id string) (string, error) {
	return s.State.GetStatus(id)
}

func (s *Service) GetNetwork(id string) (string, error) {
	return s.State.GetNetwork(id)
}
