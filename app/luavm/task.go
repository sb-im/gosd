package luavm

import "sb.im/gosd/app/model"

func (s *Service) GetTask() *model.Task {
	return s.Task
}
