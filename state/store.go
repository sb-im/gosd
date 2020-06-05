package state

import mqtt "github.com/eclipse/paho.mqtt.golang"

type State struct {
	Plan map[int]PlanState
	Node map[string]NodeState
	Mqtt mqtt.Client
}

type PlanState struct {
	id int64
}

type NodeState struct {
	Msg map[string][]byte
}
