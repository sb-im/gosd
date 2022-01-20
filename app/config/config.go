package config

type Config struct {
	Instance    string `env:"INSTANCE"`
	MqttURL     string `env:"MQTT_URL"`
	RedisURL    string `env:"REDIS_URL"`
	DatabaseURL string `env:"DATABASE_URL"`
	StorageURL  string `env:"STORAGE_URL"`
}

func DefaultConfig() *Config {
	return &Config{
		Instance:    "gosd.0",
		MqttURL:     "mqtt://admin:public@localhost:1883",
		RedisURL:    "redis://localhost:6379/1",
		DatabaseURL: "host=localhost user=postgres password=password dbname=gosd port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		StorageURL:  "file://data/storage",
	}
}
