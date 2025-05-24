package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigUnmarshal(t *testing.T) {
	yamlData := `
server:
  port: 8080
logging:
  path: "logs"
scheduler:
  max_concurrent_tasks: 5
worker:
  ping_sites:
    - "google.com"
    - "yahoo.com"
`

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	if err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server.port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Logging.Path != "logs" {
		t.Errorf("expected logging.path 'logs', got %s", cfg.Logging.Path)
	}
	if cfg.Scheduler.MaxConcurrentTasks != 5 {
		t.Errorf("expected scheduler.max_concurrent_tasks 5, got %d", cfg.Scheduler.MaxConcurrentTasks)
	}
	if len(cfg.Worker.PingSites) != 2 || cfg.Worker.PingSites[0] != "google.com" || cfg.Worker.PingSites[1] != "yahoo.com" {
		t.Errorf("unexpected worker.ping_sites: %+v", cfg.Worker.PingSites)
	}
}
