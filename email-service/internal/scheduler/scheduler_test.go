package scheduler

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault_WithEnvSet(t *testing.T) {
	os.Setenv("TEST_SCHEDULER_VAR", "custom_value")
	defer os.Unsetenv("TEST_SCHEDULER_VAR")

	result := getEnvOrDefault("TEST_SCHEDULER_VAR", "default_value")

	if result != "custom_value" {
		t.Errorf("Expected 'custom_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_WithoutEnvSet(t *testing.T) {
	os.Unsetenv("NONEXISTENT_SCHEDULER_VAR")

	result := getEnvOrDefault("NONEXISTENT_SCHEDULER_VAR", "default_value")

	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestGetEnvOrDefault_EmptyEnv(t *testing.T) {
	os.Setenv("EMPTY_SCHEDULER_VAR", "")
	defer os.Unsetenv("EMPTY_SCHEDULER_VAR")

	result := getEnvOrDefault("EMPTY_SCHEDULER_VAR", "default_value")

	// Empty string should return default
	if result != "default_value" {
		t.Errorf("Expected 'default_value' for empty env, got '%s'", result)
	}
}

func TestPostStruct_Fields(t *testing.T) {
	post := Post{
		Title: "Test Post",
		Slug:  "test-post",
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got '%s'", post.Title)
	}

	if post.Slug != "test-post" {
		t.Errorf("Expected slug 'test-post', got '%s'", post.Slug)
	}
}

func TestPostStruct_NilExcerpt(t *testing.T) {
	post := Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Excerpt: nil,
	}

	if post.Excerpt != nil {
		t.Error("Expected Excerpt to be nil")
	}
}

func TestPostStruct_WithExcerpt(t *testing.T) {
	excerpt := "This is a test excerpt"
	post := Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Excerpt: &excerpt,
	}

	if post.Excerpt == nil {
		t.Fatal("Expected Excerpt to not be nil")
	}

	if *post.Excerpt != excerpt {
		t.Errorf("Expected excerpt '%s', got '%s'", excerpt, *post.Excerpt)
	}
}

func TestSubscriptionStruct_Fields(t *testing.T) {
	sub := Subscription{
		Email: "test@example.com",
	}

	if sub.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", sub.Email)
	}
}

func TestScheduler_Close_NilDB(t *testing.T) {
	// Test that Close handles nil db gracefully
	scheduler := &Scheduler{
		db: nil,
	}

	// Should not panic
	scheduler.Close()
}
