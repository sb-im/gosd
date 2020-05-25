package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Mqtt     string `yaml:"mqtt"`
	Server   string `yaml:"server"`
	Database string `yaml:"database"`
}

func getConfig(str string) (Config, error) {
	config := Config{}
	configFile, err := ioutil.ReadFile(str)
	yaml.Unmarshal(configFile, &config)
	return config, err
}
