package mqttd

type MqttdConfig struct {
	ID      string
	Client  string
	Status  string
	Network string
	Broker  string
	Rpc     struct {
		I string
		O string
	}
	Gtran struct {
		Prefix string
	}
}

func loadMqttConfigDefault() (*MqttdConfig, error) {
	return &MqttdConfig{
		ID: "1",
		Client: "gosd-%s",
		Status: "nodes/%s/status",
		Network: "nodes/%s/network",
		// broker: "mqtt[s]://[username][:password]@host.domain[:port]"
		Broker: "mqtt://localhost:1883",
		Rpc: struct {
			I string
			O string
		} {
			I: "nodes/%s/rpc/recv",
			O: "nodes/%s/rpc/send",
		},
		Gtran: struct {
			Prefix string
		} {
			Prefix: "nodes/%s/msg/%s",
		},
	}, nil
}
