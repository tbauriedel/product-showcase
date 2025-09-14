package config

import (
	"log/slog"
	"testing"
)

// Read valid config file
func TestNewFromFile(t *testing.T) {
	c, err := NewFromFile("../../testdata/config-mysql.yaml")

	if err != nil {
		t.Fatalf("Error loading config: %s", err.Error())
	}

	if c.ListenAddr != ":1234" {
		t.Fatalf("Actual: %s, Expected: %s", c.ListenAddr, ":1234")
	}
}

// Non-existing file provided
func TestNewFromFile_NotExisting(t *testing.T) {
	_, err := NewFromFile("../../testdata/not-existing.yaml")
	if err == nil {
		t.Fatal("Expected error while providing non-existing file, got nil")
	}
}

// Invalid YAML file provided
func TestNewFromFile_InvalidYAML(t *testing.T) {
	_, err := NewFromFile("../../testdata/invalid-yaml.yaml")
	if err == nil {
		t.Fatal("Expected error while providing invalid YAML, got nil")
	}
}

// Validate log levels
func TestGetLogLevel(t *testing.T) {
	types := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"info":    slog.LevelInfo,
		"warn":    slog.LevelWarn,
		"error":   slog.LevelError,
		"default": slog.LevelInfo, // default is info
	}

	for x, y := range types {
		c := Config{LogLevel: x}

		if c.GetLogLevel() != y {
			t.Fatalf("Actual: %s, Expected: %s", c.LogLevel, y)
		}
	}
}
