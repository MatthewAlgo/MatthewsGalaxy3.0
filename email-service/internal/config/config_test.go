package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault_WithValue(t *testing.T) {
	os.Setenv("TEST_CONFIG_VAR", "custom_value")
	defer os.Unsetenv("TEST_CONFIG_VAR")

	result := GetEnvOrDefault("TEST_CONFIG_VAR", "default_value")
	if result != "custom_value" {
		t.Errorf("Expected 'custom_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_WithoutValue(t *testing.T) {
	os.Unsetenv("NONEXISTENT_CONFIG_VAR")

	result := GetEnvOrDefault("NONEXISTENT_CONFIG_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_EmptyValue(t *testing.T) {
	os.Setenv("EMPTY_CONFIG_VAR", "")
	defer os.Unsetenv("EMPTY_CONFIG_VAR")

	result := GetEnvOrDefault("EMPTY_CONFIG_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value' for empty env, got '%s'", result)
	}
}
