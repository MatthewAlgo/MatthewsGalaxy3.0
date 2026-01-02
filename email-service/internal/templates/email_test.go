package templates

import (
	"strings"
	"testing"
)

func TestNewPostEmail_GeneratesValidHTML(t *testing.T) {
	postTitle := "My Amazing Blog Post"
	postExcerpt := "This is a summary of the post content"
	postURL := "https://matthewsgalaxy.com/blog/my-amazing-post"
	unsubscribeURL := "https://matthewsgalaxy.com/api/v1/unsubscribe?email=test@example.com"

	html, err := NewPostEmail(postTitle, postExcerpt, postURL, unsubscribeURL)

	if err != nil {
		t.Fatalf("NewPostEmail returned error: %v", err)
	}

	if html == "" {
		t.Error("Expected non-empty HTML output")
	}

	// Verify key content is present
	if !strings.Contains(html, postTitle) {
		t.Errorf("HTML should contain post title: %s", postTitle)
	}

	if !strings.Contains(html, postExcerpt) {
		t.Errorf("HTML should contain post excerpt: %s", postExcerpt)
	}

	if !strings.Contains(html, postURL) {
		t.Errorf("HTML should contain post URL: %s", postURL)
	}

	if !strings.Contains(html, unsubscribeURL) {
		t.Errorf("HTML should contain unsubscribe URL: %s", unsubscribeURL)
	}
}

func TestNewPostEmail_ContainsRequiredElements(t *testing.T) {
	html, err := NewPostEmail("Test Post", "Test excerpt", "https://example.com/post", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("NewPostEmail returned error: %v", err)
	}

	requiredElements := []string{
		"<!DOCTYPE html>",
		"<html>",
		"</html>",
		"Matthew's Galaxy",
		"New Post Alert",
		"Read the Full Post",
		"Unsubscribe",
	}

	for _, element := range requiredElements {
		if !strings.Contains(html, element) {
			t.Errorf("HTML should contain: %s", element)
		}
	}
}

func TestNewPostEmail_HandlesSpecialCharacters(t *testing.T) {
	postTitle := "Test <script>alert('xss')</script> Post"
	postExcerpt := "Content with \"quotes\" and 'apostrophes'"

	html, err := NewPostEmail(postTitle, postExcerpt, "https://example.com", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("NewPostEmail returned error: %v", err)
	}

	// HTML template should escape special characters
	// The template engine escapes < and > in content
	if strings.Contains(html, "<script>") {
		t.Error("HTML should escape script tags for security")
	}
}

func TestWelcomeEmail_GeneratesValidHTML(t *testing.T) {
	userName := "John Doe"
	unsubscribeURL := "https://matthewsgalaxy.com/api/v1/unsubscribe?email=john@example.com"

	html, err := WelcomeEmail(userName, unsubscribeURL)

	if err != nil {
		t.Fatalf("WelcomeEmail returned error: %v", err)
	}

	if html == "" {
		t.Error("Expected non-empty HTML output")
	}

	// Verify key content is present
	if !strings.Contains(html, userName) {
		t.Errorf("HTML should contain user name: %s", userName)
	}

	if !strings.Contains(html, unsubscribeURL) {
		t.Errorf("HTML should contain unsubscribe URL: %s", unsubscribeURL)
	}
}

func TestWelcomeEmail_ContainsRequiredElements(t *testing.T) {
	html, err := WelcomeEmail("Test User", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("WelcomeEmail returned error: %v", err)
	}

	requiredElements := []string{
		"<!DOCTYPE html>",
		"<html>",
		"</html>",
		"Matthew's Galaxy",
		"Welcome aboard",
		"successfully subscribed",
		"Unsubscribe",
		"Cloud Architecture",
		"Software Engineering",
	}

	for _, element := range requiredElements {
		if !strings.Contains(html, element) {
			t.Errorf("HTML should contain: %s", element)
		}
	}
}

func TestWelcomeEmail_HandlesEmptyUserName(t *testing.T) {
	html, err := WelcomeEmail("", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("WelcomeEmail returned error: %v", err)
	}

	// Should still generate valid HTML even with empty username
	if !strings.Contains(html, "Welcome aboard") {
		t.Error("HTML should still contain welcome message")
	}
}

func TestNewPostEmail_ContainsSpaceTheme(t *testing.T) {
	html, err := NewPostEmail("Test", "Test", "https://example.com", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("NewPostEmail returned error: %v", err)
	}

	// Check for space-themed elements
	spaceElements := []string{
		"✨",
		"🌟",
		"🚀",
	}

	foundAny := false
	for _, element := range spaceElements {
		if strings.Contains(html, element) {
			foundAny = true
			break
		}
	}

	if !foundAny {
		t.Error("HTML should contain space-themed emoji elements")
	}
}

func TestWelcomeEmail_ContainsSpaceTheme(t *testing.T) {
	html, err := WelcomeEmail("Test User", "https://example.com/unsub")

	if err != nil {
		t.Fatalf("WelcomeEmail returned error: %v", err)
	}

	// Check for space-themed motto
	if !strings.Contains(html, "Ad astra per aspera") {
		t.Error("HTML should contain the space-themed motto")
	}
}
