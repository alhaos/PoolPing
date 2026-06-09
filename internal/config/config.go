// Package config provides configuration loading for PoolPing.
// Settings are read from a YAML file and can be overridden via environment variables.
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds all runtime configuration for PoolPing.
type Config struct {
	// Workers is the number of parallel ping goroutines.
	// Env: POOLPING_WORKERS, default: 5.
	Workers int `yaml:"workers" env:"POOLPING_WORKERS" env-default:"5"`

	// TimeoutMs is the ICMP reply timeout in milliseconds.
	// Env: POOLPING_TIMEOUT_MS, default: 3000.
	TimeoutMs int `yaml:"timeoutMs" env:"POOLPING_TIMEOUT_MS" env-default:"3000"`

	// Hosts is the list of hostnames or IP addresses to ping.
	// Env: POOLPING_HOSTS (comma-separated).
	Hosts []string `yaml:"hosts" env:"POOLPING_HOSTS" env-separator:","`
}

// Load reads configuration from the YAML file at path
// and applies any environment variable overrides.
func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
