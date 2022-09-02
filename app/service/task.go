package service

import (
	"context"
	"encoding/json"

	"sb.im/gosd/app/logger"
	"sb.im/gosd/app/model"
)

func (s *Service) TaskRun(ctx context.Context, task *model.Task) error {
	files := make(map[string]string)
	extra := make(map[string]string)

	json.Unmarshal(task.Files, &files)
	json.Unmarshal(task.Extra, &extra)

	logger.WithContext(ctx).Infof("%+v\t%v\t%v", task, files, extra)
	return s.worker.AddTask(ctx, task)
}

func (s *Service) TaskKill(taskId string) error {
	return s.worker.Kill(taskId)
}
