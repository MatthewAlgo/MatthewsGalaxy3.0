package sender

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault_WithEnvSet(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "custom_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := getEnvOrDefault("TEST_ENV_VAR", "default_value")

	if result != "custom_value" {
		t.Errorf("Expected 'custom_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_WithoutEnvSet(t *testing.T) {
	os.Unsetenv("NONEXISTENT_VAR")

	result := getEnvOrDefault("NONEXISTENT_VAR", "default_value")

	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_EmptyEnv(t *testing.T) {
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result := getEnvOrDefault("EMPTY_VAR", "default_value")

	// Empty string should return default
	if result != "default_value" {
		t.Errorf("Expected 'default_value' for empty env, got '%s'", result)
	}
}

func TestNewEmailSender_DefaultConfig(t *testing.T) {
	// Clear environment
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USER")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("FROM_EMAIL")

	sender := NewEmailSender()

	if sender == nil {
		t.Fatal("NewEmailSender returned nil")
	}

	if sender.host != "smtp.gmail.com" {
		t.Errorf("Expected default host 'smtp.gmail.com', got '%s'", sender.host)
	}

	if sender.port != 587 {
		t.Errorf("Expected default port 587, got %d", sender.port)
	}

	if sender.from != "noreply@matthewsgalaxy.com" {
		t.Errorf("Expected default from email, got '%s'", sender.from)
	}
}

func TestNewEmailSender_CustomConfig(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.custom.com")
	os.Setenv("SMTP_PORT", "465")
	os.Setenv("SMTP_USER", "user@custom.com")
	os.Setenv("SMTP_PASSWORD", "secret123")
	os.Setenv("FROM_EMAIL", "no-reply@custom.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USER")
		os.Unsetenv("SMTP_PASSWORD")
		os.Unsetenv("FROM_EMAIL")
	}()

	sender := NewEmailSender()

	if sender.host != "smtp.custom.com" {
		t.Errorf("Expected custom host, got '%s'", sender.host)
	}

	if sender.port != 465 {
		t.Errorf("Expected port 465, got %d", sender.port)
	}

	if sender.username != "user@custom.com" {
		t.Errorf("Expected custom username, got '%s'", sender.username)
	}

	if sender.password != "secret123" {
		t.Errorf("Expected custom password, got '%s'", sender.password)
	}

	if sender.from != "no-reply@custom.com" {
		t.Errorf("Expected custom from email, got '%s'", sender.from)
	}
}

func TestNewEmailSender_InvalidPort(t *testing.T) {
	os.Setenv("SMTP_PORT", "not-a-number")
	defer os.Unsetenv("SMTP_PORT")

	sender := NewEmailSender()

	// Invalid port should default to 587
	if sender.port != 587 {
		t.Errorf("Expected default port 587 for invalid port string, got %d", sender.port)
	}
}

func TestSendBulkEmails_EmptyRecipients(t *testing.T) {
	sender := NewEmailSender()
	recipients := []string{}

	errors := sender.SendBulkEmails(recipients, "Test Subject", func(email string) string {
		return "<html>Test</html>"
	})

	// No recipients means no errors
	if len(errors) != 0 {
		t.Errorf("Expected no errors for empty recipients, got %d", len(errors))
	}
}
