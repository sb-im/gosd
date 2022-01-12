package service

import (
	"encoding/json"
	"fmt"

	"sb.im/gosd/app/model"

	"sb.im/gosd/luavm"
)

func (s *Service) TaskRun(task *model.Task) error {
	fmt.Println(task)

	files := make(map[string]string)
	extra := make(map[string]string)

	json.Unmarshal(task.Files, &files)
	json.Unmarshal(task.Extra, &extra)

	fmt.Println(files)
	fmt.Println(extra)

	// TODO: join this worker
	s.worker.Queue <- &luavm.Task{

		// TODO:
		ID:     int64(task.ID),
		NodeID: task.NodeID,
		PlanID: int64(task.ID),

		Files: files,
		Extra: extra,

		// TODO:
		Script: []byte{},
	}
	return nil
}
