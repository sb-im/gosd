package service

import (
	"fmt"

	"sb.im/gosd/app/model"
)

func (s *Service) AddSchedule(task *model.Schedule) error {
	_, err := s.cron.AddFunc(task.Cron, func() {
		fmt.Println("Cron running:", task.Name)
		s.JSON.Call(task.Target, []byte(task.Params))
	})
	return err
}
