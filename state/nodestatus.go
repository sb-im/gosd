package state

import (
	"encoding/json"
)

const (
	online  = 0
	offline = 1
	neterr  = 2
)

type NodeStatus struct {
	Code   int    `json:"code"`
	Status Status `json:"status"`
}

// [Reference] https://gitlab.com/sbim/superdock/cloud/ncp/-/blob/master/config.go#L35
type Status struct {
	LinkId      int  `json:"link_id"`
	Position_ok bool `json:"position_ok"`
	Position
}

type Position struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
	Alt string `json:"alt"`
}

func (n *NodeStatus) isConnect() bool {
	connect := false
	if n.Code == 0 {
		connect = true
	}
	return connect
}

func (n *NodeStatus) GetID(str string) int {
	return n.Status.LinkId
}

func (n *NodeStatus) SetStatus(raw []byte) error {
	err := json.Unmarshal(raw, n)
	if err != nil {
		return err
	}
	return nil
}
