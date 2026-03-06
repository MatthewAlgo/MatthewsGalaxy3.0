package scheduler

import (
	"os"
	"testing"

	"github.com/matthewsgalaxy/email-service/internal/config"
)

func TestGetEnvOrDefault_WithValue(t *testing.T) {
	os.Setenv("TEST_SCHEDULER_VAR", "test_value")
	defer os.Unsetenv("TEST_SCHEDULER_VAR")

	result := config.GetEnvOrDefault("TEST_SCHEDULER_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_SCHEDULER_VAR")

	result := config.GetEnvOrDefault("NONEXISTENT_SCHEDULER_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_EmptyValue(t *testing.T) {
	os.Setenv("EMPTY_SCHEDULER_VAR", "")
	defer os.Unsetenv("EMPTY_SCHEDULER_VAR")

	result := config.GetEnvOrDefault("EMPTY_SCHEDULER_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value' for empty env, got '%s'", result)
	}
}
