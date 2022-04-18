package service

import (
	"context"
	"fmt"

	"sb.im/gosd/app/helper"
	"sb.im/gosd/app/model"
)

const (
	mqttAuthReqPrefix = "mqtt_user:"
	mqttAuthAclPrefix = "mqtt_acl:"

	mqttAuthTeam = "user.%d"
	mqttAuthNode = "node.%s"

	mqttTopicNode = "nodes/%s/#"
	mqttTopicTask = "tasks/%d/#"
)

// Emqx auth acl plugin
// access: Allowed operations: subscribe (1), publish (2), both subscribe and publish (3)
const (
	mqttAuthAccessSubscribe = 1
	mqttAuthAccessPublish   = 2
	mqttAuthAccessPubSub    = 3
)

func (s *Service) MqttAuthReqTeam(teamID uint) (username string, password string, err error) {
	username = fmt.Sprintf(mqttAuthTeam, teamID)
	password = helper.GenSecret(8)
	err = s.rdb.HSet(context.Background(), mqttAuthReqPrefix+username, map[string]interface{}{
		"password": password,
	}).Err()
	return
}

func (s *Service) MqttAuthAclTeam(teamID uint) error {
	username := fmt.Sprintf(mqttAuthTeam, teamID)

	var nodes []model.Node
	if err := s.orm.Find(&nodes, "team_id = ?", teamID).Error; err != nil {
		return err
	}

	acl := make(map[string]interface{})
	for _, node := range nodes {
		acl[fmt.Sprintf(mqttTopicNode, node.ID)] = mqttAuthAccessPubSub
	}

	var tasks []model.Task
	if err := s.orm.Find(&tasks, "team_id = ?", teamID).Error; err != nil {
		return err
	}
	for _, task := range tasks {
		acl[fmt.Sprintf(mqttTopicTask, task.ID)] = mqttAuthAccessPubSub
	}

	return s.rdb.HSet(context.Background(), mqttAuthAclPrefix+username, acl).Err()
}

func (s *Service) MqttAuthReqNode(nodeID string) error {
	username := fmt.Sprintf(mqttAuthNode, nodeID)

	var node model.Node
	if err := s.orm.Take(&node, "id", nodeID).Error; err != nil {
		return err
	}

	return s.rdb.HSet(context.Background(), mqttAuthReqPrefix+username, map[string]interface{}{
		"password": node.Secret,
	}).Err()
}

func (s *Service) MqttAuthAclNode(nodeID string) error {
	username := fmt.Sprintf(mqttAuthNode, nodeID)

	return s.rdb.HSet(context.Background(), mqttAuthAclPrefix+username, map[string]interface{}{
		fmt.Sprintf(mqttTopicNode, nodeID): mqttAuthAccessPubSub,
	}).Err()
}

func (s *Service) MqttAuthNodeSync() error {
	var nodes []model.Node
	if err := s.orm.Find(&nodes).Error; err != nil {
		return err
	}
	for _, v := range nodes {
		if err := s.MqttAuthReqNode(v.ID); err != nil {
			return err
		}
		if err := s.MqttAuthAclNode(v.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) MqttAuthTeamSync() error {
	var teams []model.Team
	if err := s.orm.Find(&teams).Error; err != nil {
		return err
	}
	for _, team := range teams {
		s.MqttAuthAclTeam(team.ID)
	}
	return nil
}
