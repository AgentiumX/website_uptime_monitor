package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the root configuration structure.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Agent    AgentConfig    `yaml:"agent"`
	Database DatabaseConfig `yaml:"database"`
	TSDB     TSDBConfig     `yaml:"tsdb"`
	Alert    AlertConfig    `yaml:"alert"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port          int    `yaml:"port"`
	JWTSecret     string `yaml:"jwt_secret"`
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
}

// AgentConfig holds agent communication settings.
type AgentConfig struct {
	SharedSecret     string `yaml:"shared_secret"`
	HeartbeatTimeout int    `yaml:"heartbeat_timeout"`
}

// DatabaseConfig holds SQLite database settings.
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// TSDBConfig holds Prometheus TSDB settings.
type TSDBConfig struct {
	Path         string `yaml:"path"`
	RetentionRaw string `yaml:"retention_raw"`
	Retention5m  string `yaml:"retention_5m"`
	Retention1h  string `yaml:"retention_1h"`
}

// AlertConfig holds alert-related settings.
type AlertConfig struct {
	HistoryRetention int `yaml:"history_retention"`
}

// Load reads and parses a YAML configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
