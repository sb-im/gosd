package service

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"sb.im/gosd/app/model"
)

func (s *Service) TaskRun(task *model.Task) error {
	log.Println(task)

	files := make(map[string]string)
	extra := make(map[string]string)

	json.Unmarshal(task.Files, &files)
	json.Unmarshal(task.Extra, &extra)

	log.Println("files: ", files)
	log.Println("extra: ", extra)

	return s.worker.AddTask(task)
}

func (s *Service) TaskKill(taskId string) error {
	return s.worker.Kill(taskId)
}
