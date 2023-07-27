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
	// Auth
	Secret    string `env:"SECRET"`
	ApiKey    string `env:"API_KEY"`
	BasicAuth bool   `env:"BASIC_AUTH"`

	// Service
	MqttURL     string `env:"MQTT_URL"`
	RedisURL    string `env:"REDIS_URL"`
	StorageURL  string `env:"STORAGE_URL"`
	DatabaseURL string `env:"DATABASE_URL"`

	// Public
	BaseURL    string `env:"BASE_URL"`
	ClientURL  string `env:"CLIENT_URL"`
	ApiMqttWs  string `env:"API_MQTT_WS"`
	ListenAddr string `env:"LISTEN_ADDR"`

	// Feature Flags
	ResetMode   bool   `env:"RESET_MODE"`
	Schedule    bool   `env:"SCHEDULE"`
	EmqxAuth    bool   `env:"EMQX_AUTH"`
	LuaFilePath string `env:"LUA_FILE"`
	LogLevel    string `env:"LOG_LEVEL"`

	// Custom
	Instance string `env:"INSTANCE"`
	Language string `env:"LANGUAGE"`
	Timezone string `env:"TIMEZONE"`

	// Developer
	Debug      bool `env:"DEBUG"`
	DemoMode   bool `env:"DEMO_MODE"`
	SingleUser bool `env:"SINGLE_USER"`
}

var opts = DefaultConfig()

func DefaultConfig() *Config {
	return &Config{
		Secret:    "falling-cats-and-dogs",
		ApiKey:    "the-elephant-in-the-room",
		BasicAuth: true,

		MqttURL:     "mqtt://admin:public@localhost:1883",
		RedisURL:    "redis://localhost:6379/1",
		StorageURL:  "data/storage",
		DatabaseURL: "postgres://postgres:password@localhost:5432/gosd?sslmode=disable&TimeZone=Asia/Shanghai",

		BaseURL:    "http://localhost:8000/gosd/api/v3",
		ClientURL:  "http://localhost:8000/gosd/api/v3",
		ApiMqttWs:  "ws://localhost:8083/mqtt",
		ListenAddr: "0.0.0.0:8000",

		ResetMode:   true,
		Schedule:    true,
		EmqxAuth:    false,
		LuaFilePath: "default.lua",

		Instance: "gosd",
		Language: "en_US",
		Timezone: "Asia/Shanghai",
		LogLevel: "info",

		Debug:      true,
		DemoMode:   false,
		SingleUser: true,
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
