package luavm

import "sb.im/gosd/app/model"

func (s *Service) GetNode(id string) *model.Node {
	for _, node := range s.nodes {
		if node.ID == id {
			return &node
		}
	}
	return nil
}
