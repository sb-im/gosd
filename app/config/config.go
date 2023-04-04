package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	loadEnvConfig()
}

type Config struct {
	ListenAddr  string `env:"LISTEN_ADDR"`
	Instance    string `env:"INSTANCE"`
	BaseURL     string `env:"BASE_URL"`
	MqttURL     string `env:"MQTT_URL"`
	RedisURL    string `env:"REDIS_URL"`
	ClientURL   string `env:"CLIENT_URL"`
	DatabaseURL string `env:"DATABASE_URL"`
	StorageURL  string `env:"STORAGE_URL"`
	LuaFilePath string `env:"LUA_FILE"`
	Debug       bool   `env:"DEBUG"`
	SingleUser  bool   `env:"SINGLE_USER"`
	BasicAuth   bool   `env:"BASIC_AUTH"`
	DemoMode    bool   `env:"DEMO_MODE"`
	Schedule    bool   `env:"SCHEDULE"`
	EmqxAuth    bool   `env:"EMQX_AUTH"`
	Language    string `env:"LANGUAGE"`
	Timezone    string `env:"TIMEZONE"`
	ApiKey      string `env:"API_KEY"`
	ApiMqtt     string `env:"API_MQTT"`
	ApiMqttWs   string `env:"API_MQTT_WS"`
	Secret      string `env:"SECRET"`
}

var opts = DefaultConfig()

func DefaultConfig() *Config {
	return &Config{
		ListenAddr:  "0.0.0.0:8000",
		Instance:    "gosd.0",
		BaseURL:     "http://localhost:8000/gosd/api/v3",
		MqttURL:     "mqtt://admin:public@localhost:1883",
		RedisURL:    "redis://localhost:6379/1",
		ClientURL:   "http://localhost:8000/gosd/api/v3",
		DatabaseURL: "postgres://postgres:password@localhost:5432/gosd?sslmode=disable&TimeZone=Asia/Shanghai",
		StorageURL:  "data/storage",
		LuaFilePath: "default.lua",
		Debug:       true,
		SingleUser:  true,
		BasicAuth:   true,
		DemoMode:    false,
		Schedule:    true,
		EmqxAuth:    false,
		Language:    "en_US",
		Timezone:    "Asia/Shanghai",
		ApiMqtt:     "mqtt://localhost:1883",
		ApiMqttWs:   "ws://localhost:8083/mqtt",
	}
}

func loadDotEnvConfig() error {
	return godotenv.Load()
}

func loadEnvConfig() error {
	return env.Parse(opts)
}

func Parse(args ...string) *Config {
	if err := loadDotEnvConfig(); err != nil {
		log.Debug(err)
	}
	if err := loadEnvConfig(); err != nil {
		log.Info(err)
	}

	log.Debugf("%+v\n", *opts)

	return opts
}

func Opts() *Config {
	return opts
}
