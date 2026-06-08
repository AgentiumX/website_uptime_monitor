package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all agent configuration loaded from agent.yaml.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Agent  AgentConfig  `yaml:"agent"`
	Probe  ProbeConfig  `yaml:"probe"`
	Report ReportConfig `yaml:"report"`
}

// ServerConfig holds the server connection settings.
type ServerConfig struct {
	URL          string `yaml:"url"`
	SharedSecret string `yaml:"shared_secret"`
}

// AgentConfig holds the agent identity settings.
type AgentConfig struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

// ProbeConfig holds probe engine settings.
type ProbeConfig struct {
	Timeout        int `yaml:"timeout"`
	MaxConcurrent  int `yaml:"max_concurrent"`
}

// ReportConfig holds result reporting settings.
type ReportConfig struct {
	Interval    int `yaml:"interval"`
	BatchSize   int `yaml:"batch_size"`
	RetryMaxAge int `yaml:"retry_max_age"`
}

// LoadConfig reads a YAML config file and returns a Config with defaults applied.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// Apply defaults
	if cfg.Probe.Timeout <= 0 {
		cfg.Probe.Timeout = 30
	}
	if cfg.Probe.MaxConcurrent <= 0 {
		cfg.Probe.MaxConcurrent = 50
	}
	if cfg.Report.Interval <= 0 {
		cfg.Report.Interval = 30
	}
	if cfg.Report.BatchSize <= 0 {
		cfg.Report.BatchSize = 100
	}
	if cfg.Report.RetryMaxAge <= 0 {
		cfg.Report.RetryMaxAge = 600
	}

	return cfg, nil
}
