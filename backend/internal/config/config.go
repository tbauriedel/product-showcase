package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	defaultLogLevel   = "info"
	defaultListenAddr = ":4890"
)

type Config struct {
	LogLevel   string   `yaml:"log_level,omitempty"`
	ListenAddr string   `yaml:"listener,omitempty"`
	Database   Database `yaml:"database,omitempty"`
}

type Database struct {
	Type    string `yaml:"type,omitempty"`
	Address string `yaml:"address,omitempty"`
	Port    int    `yaml:"port,omitempty"`
	User    string `yaml:"user,omitempty"`
	Pass    string `yaml:"pass,omitempty"`
	DbName  string `yaml:"name,omitempty"`
}

// NewFromFile reads the configuration from a YAML file and returns a Config instance.
func NewFromFile(filename string) (Config, error) {
	// Set default config values
	conf := Config{
		LogLevel:   defaultLogLevel,
		ListenAddr: defaultListenAddr,
	}

	// Open provided config file
	f, err := os.Open(filename)
	if err != nil {
		return conf, fmt.Errorf("error reading config file: %w", err)
	}

	// Read and parse config file
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&conf); err != nil {
		return conf, fmt.Errorf("error parsing config file: %w", err)
	}

	return conf, nil
}

// GetLogLevel returns the loglevel
func (c *Config) GetLogLevel() slog.Level {
	// Set log level given via flag
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
