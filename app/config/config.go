package config

import (
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func init() {
	loadEnvConfig()
}

type Config struct {
	Instance    string `env:"INSTANCE"`
	BaseURL     string `env:"BASE_URL"`
	MqttURL     string `env:"MQTT_URL"`
	RedisURL    string `env:"REDIS_URL"`
	DatabaseURL string `env:"DATABASE_URL"`
	StorageURL  string `env:"STORAGE_URL"`
	LuaFilePath string `env:"LUA_FILE"`
	Debug       bool   `env:"DEBUG" yaml:"debug"`
	SingleUser  bool   `env:"SINGLE_USER" yaml:"single_user"`
	BasicAuth   bool   `env:"BASIC_AUTH" yaml:"basic_auth"`
	Language    string `env:"LANGUAGE" yaml:"language"`
	Timezone    string `env:"TIMEZONE" yaml:"timezone"`
	ApiKey      string `env:"API_KEY" yaml:"api_key"`
	ApiMqtt     string `env:"API_MQTT"`
	ApiMqttWs   string `env:"API_MQTT_WS"`
	Secret      string `env:"SECRET" yaml:"secret"`
}

var opts = DefaultConfig()

func DefaultConfig() *Config {
	return &Config{
		Instance:    "gosd.0",
		BaseURL:     "http://localhost:8000/gosd/api/v3",
		MqttURL:     "mqtt://admin:public@localhost:1883",
		RedisURL:    "redis://localhost:6379/1",
		DatabaseURL: "host=localhost user=postgres password=password dbname=gosd port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		StorageURL:  "file://data/storage",
		LuaFilePath: "default.lua",
		Debug:       true,
		SingleUser:  true,
		BasicAuth:   true,
		Language:    "en_US",
		Timezone:    "Asia/Shanghai",
		ApiMqtt:     "mqtt://localhost:1883",
		ApiMqttWs:   "ws://localhost:8083/mqtt",
	}
}

func loadYamlConfig(str string) error {
	configFile, err := ioutil.ReadFile(str)
	if err != nil {
		return err
	} else {
		err = yaml.Unmarshal(configFile, &opts)
		return err
	}
}

func loadDotEnvConfig() error {
	return godotenv.Load()
}

func loadEnvConfig() error {
	return env.Parse(opts)
}

func Parse(args ...string) *Config {
	if len(args) >= 1 {
		if err := loadYamlConfig(args[0]); err != nil {
			log.Debug(err)
		}
	}
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
