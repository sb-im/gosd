package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"sb.im/gosd/app/model"
)

const (
	mqttAuthUserPrefix = "mqtt_user:"
	mqttAuthACLPrefix  = "mqtt_acl:"
)

// Emqx auth acl plugin
// access: Allowed operations: subscribe (1), publish (2), both subscribe and publish (3)
const (
	mqttAuthAccessSubscribe = 1
	mqttAuthAccessPublish   = 2
	mqttAuthAccessPubSub    = 3
)

func genToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *Service) MqttAuthUser(user string) string {
	password, _ := genToken(8)
	s.rdb.HSet(context.Background(), mqttAuthUserPrefix+user, map[string]interface{}{
		"password": password,
	})
	return password
}

func (s *Service) MqttAuthACL(teamID uint, user string) string {
	var nodes []model.Node
	s.orm.Find(&nodes, "team_id = ?", teamID)

	acl := make(map[string]interface{})
	for _, node := range nodes {
		acl["nodes/"+strconv.Itoa(int(node.ID))+"/#"] = mqttAuthAccessPubSub
	}

	s.rdb.HSet(context.Background(), mqttAuthACLPrefix+user, acl)
	return ""
}

func (s *Service) MqttAuthNodeUser(nodeID string) error {
	var node model.Node
	if err := s.orm.Take(&node, nodeID).Error; err != nil {
		return err
	}

	return s.rdb.HSet(context.Background(), mqttAuthUserPrefix+nodeID, map[string]interface{}{
		"password": node.Secret,
	}).Err()
}

func (s *Service) MqttAuthNodeACL(nodeID string) error {
	return s.rdb.HSet(context.Background(), mqttAuthACLPrefix+nodeID, map[string]interface{}{
		"nodes/" + nodeID + "/#": mqttAuthAccessPubSub,
	}).Err()
}

func (s *Service) MqttAuthSync() error {
	var nodes []model.Node
	if err := s.orm.Find(&nodes).Error; err != nil {
		return err
	}
	for _, v := range nodes {
		strId := strconv.Itoa(int(v.ID))
		if err := s.MqttAuthNodeUser(strId); err != nil {
			return err
		}
		if err := s.MqttAuthNodeACL(strId); err != nil {
			return err
		}
	}
	return nil
}
