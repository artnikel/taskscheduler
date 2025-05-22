package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LoggingConfig struct {
	Path string `yaml:"path"`
}

type SchedulerConfig struct {
	MaxConcurrentTasks int `yaml:"max_concurrent_tasks"`
}

type WorkerConfig struct {
	PingSites []string `yaml:"ping_sites"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Worker    WorkerConfig    `yaml:"worker"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
