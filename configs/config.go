package configs

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Server      string        `yaml:"server"`
	Interval    time.Duration `yaml:"interval"`
	MaxRequests int           `yaml:"max_requests"`
}

func LoadConfig(name string) (Config, error) {
	yamlFile, err := os.ReadFile(name)
	if err != nil {
		return Config{}, err
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
