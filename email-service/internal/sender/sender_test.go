package sender

import (
	"testing"

	"os"

	"github.com/matthewsgalaxy/email-service/internal/config"
)

func TestGetEnvOrDefault_WithValue(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := config.GetEnvOrDefault("TEST_ENV_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_VAR")

	result := config.GetEnvOrDefault("NONEXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_EmptyValue(t *testing.T) {
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result := config.GetEnvOrDefault("EMPTY_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value' for empty env, got '%s'", result)
	}
}

func TestNewEmailSender_DefaultPort(t *testing.T) {
	os.Unsetenv("SMTP_PORT")
	s := NewEmailSender()
	if s.port != 587 {
		t.Errorf("Expected default port 587, got %d", s.port)
	}
}

func TestNewEmailSender_CustomPort(t *testing.T) {
	os.Setenv("SMTP_PORT", "465")
	defer os.Unsetenv("SMTP_PORT")

	s := NewEmailSender()
	if s.port != 465 {
		t.Errorf("Expected port 465, got %d", s.port)
	}
}

func TestNewEmailSender_CustomHost(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.custom.com")
	defer os.Unsetenv("SMTP_HOST")

	s := NewEmailSender()
	if s.host != "smtp.custom.com" {
		t.Errorf("Expected host 'smtp.custom.com', got '%s'", s.host)
	}
}
