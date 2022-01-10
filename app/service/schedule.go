package service

import (
	"fmt"

	"sb.im/gosd/app/model"
)

func (s Service) StartSchedule() {
	schedules := []model.Schedule{}
	s.orm.Find(&schedules)

	for _, i := range schedules {
		s.cron.AddFunc(i.Cron, func() {
			s.JSON.Call(i.Method, []byte(i.Params))
		})
	}

	s.cron.Start()
}

func (s *Service) AddSchedule(task *model.Schedule) error {
	_, err := s.cron.AddFunc(task.Cron, func() {
		fmt.Println("Cron running:", task.Name)
		s.JSON.Call(task.Method, []byte(task.Params))
	})
	return err
}
