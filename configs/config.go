package configs

import "time"

type Config struct {
	Server      string        `yaml:"server"`
	WindowSize  time.Duration `yaml:"window_size"`
	MaxRequests int           `yaml:"max_requests"`
}
