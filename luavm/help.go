package luavm

import (
	"time"
)

// Reference: https://golang.org/pkg/time/#ParseDuration
func (s *Service) Sleep(str string) error {
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	time.Sleep(duration)
	return nil
}
