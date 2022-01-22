package service

import (
	"encoding/json"
	"fmt"

	"sb.im/gosd/app/model"
)

func (s Service) TaskRun(task *model.Task) error {
	fmt.Println(task)

	files := make(map[string]string)
	extra := make(map[string]string)

	json.Unmarshal(task.Files, &files)
	json.Unmarshal(task.Extra, &extra)

	fmt.Println("files: ", files)
	fmt.Println("extra: ", extra)

	// TODO: join this worker
	s.worker.RunTask(task)
	return nil
}
