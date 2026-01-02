package handlers

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "title with special characters",
			title:    "Hello, World! How are you?",
			expected: "hello-world-how-are-you",
		},
		{
			name:     "title with numbers",
			title:    "Top 10 Tips for 2026",
			expected: "top-10-tips-for-2026",
		},
		{
			name:     "title with multiple spaces",
			title:    "Hello    World",
			expected: "hello-world",
		},
		{
			name:     "title with leading/trailing spaces",
			title:    "  Hello World  ",
			expected: "hello-world",
		},
		{
			name:     "title with unicode",
			title:    "Café & Restaurant",
			expected: "caf-restaurant",
		},
		{
			name:     "empty title",
			title:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			title:    "!@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSlug(tt.title)
			if result != tt.expected {
				t.Errorf("generateSlug(%q) = %q, want %q", tt.title, result, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "simple number",
			input:    "123",
			expected: 123,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "number with letters",
			input:    "abc123",
			expected: 123,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "only letters",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "mixed format",
			input:    "12ab34",
			expected: 1234,
		},
		{
			name:     "large number",
			input:    "999999",
			expected: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
