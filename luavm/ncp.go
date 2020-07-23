package luavm

import (
	"fmt"
	"os"
)

const (
	blob_url = "api/v1/blobs/%s"
)

func (s *LService) FileUrl(key string) string {
	return fmt.Sprintf(os.Getenv("BASE_URL")+blob_url, s.Task.Attach[key])
}
