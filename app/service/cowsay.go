package service

import (
	"fmt"
)

func (s *Service) CowSay(content string) error {
	fmt.Println(content)
	return nil
}
