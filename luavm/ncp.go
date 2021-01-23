package luavm

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"sb.im/gosd/model"
)

const (
	blob_url = "api/v1/blobs/%s"
)

func (s *Service) FileUrl(key string) string {
	blobID, ok := s.Task.Attach[key]
	if !ok {
		id, err := strconv.ParseInt(s.Task.PlanID, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		plan, err := s.Store.PlanByID(id)
		if err != nil {
			fmt.Println(err)
		}

		blob := model.NewBlob("", bytes.NewReader(nil))
		s.Store.CreateBlob(blob)

		blobID = strconv.FormatInt(blob.ID, 10)
		plan.Attachments[key] = blobID

		err = s.Store.UpdatePlan(plan)
		if err != nil {
			fmt.Println(err)
		}
	}
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, blobID)
}

func (s *Service) LogFileUrl(key string) string {
	id, err := strconv.ParseInt(s.Task.ID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	plan, err := s.Store.PlanLogByID(id)
	if err != nil {
		fmt.Println(err)
	}

	blob := model.NewBlob("", bytes.NewReader(nil))
	s.Store.CreateBlob(blob)

	blobID := strconv.FormatInt(blob.ID, 10)
	plan.Attachments[key] = blobID

	err = s.Store.UpdatePlanLog(plan)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, blobID)
}
