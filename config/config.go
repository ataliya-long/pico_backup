package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DBNode stores connection defaults for one database type.
type DBNode struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Config holds PostgreSQL connection defaults.
type Config struct {
	Postgres DBNode `yaml:"postgres"`
	BackupDir string `yaml:"backup_dir"` // output directory
}

// Default returns hard-coded defaults.
func Default() *Config {
	return &Config{
		BackupDir: ".",
		Postgres: DBNode{Host: "192.168.23.129", Port: 5432, User: "postgres"},
	}
}

// Load reads a YAML file and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, nil
}
