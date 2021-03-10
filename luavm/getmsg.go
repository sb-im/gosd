package luavm

func (s *Service) GetMsg(id, msg string) (string, error) {
	raw, err := s.State.GetNodeMsg(id, msg)
	return string(raw), err
}

func (s *Service) GetStatus(id string) (string, error) {
	return s.State.GetNode(id, "status")
}

func (s *Service) GetNetwork(id string) (string, error) {
	return s.State.GetNode(id, "network")
}
