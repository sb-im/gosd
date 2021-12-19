package service

import (
	"context"
	"fmt"
	//"sb.im/gosd/app/model"
)

const (
	mqttAuthUserPrefix = "mqtt_user:"
	mqttAuthACLPrefix  = "mqtt_acl:"
)

func (s *Service) MqttAuthUser(user string) string {
	// TODO: need random generate password
	password := "xxx"
	fmt.Println(user, password)
	s.rdb.HSet(context.Background(), mqttAuthUserPrefix+user, map[string]interface{}{
		"password": password,
	})
	return ""
}

func (s *Service) MqttAuthACL(user string) string {
	fmt.Println(user)
	// TODO:
	s.rdb.HSet(context.Background(), mqttAuthACLPrefix+user, map[string]interface{}{
		"nodes":   true,
		"nodes/1": true,
		"nodes/#": true,
	})

	return ""
}
