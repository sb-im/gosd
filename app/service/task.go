package service

import (
	"fmt"

	"sb.im/gosd/app/model"
	"sb.im/gosd/luavm"
)

func (s *Service) PlanTask(planID string) {
	task := &model.Task{}
	//task := &model.Plan{}
	s.orm.Find(task, planID)
	//s.orm.Table("plans").Select("name", "extra").Where("id = ?", planID).Scan(&task)
	fmt.Println(task)
}

func (s *Service) TaskRun(task *luavm.Task) {
	s.worker.Queue <- task
}
