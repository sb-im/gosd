package service

import (
	"fmt"

	"sb.im/gosd/app/model"
)

func (s *Service) TaskRun(task *model.Task) error {
	fmt.Println(task)

	// TODO: join this worker
	//s.worker.Queue <- task
	return nil
}
