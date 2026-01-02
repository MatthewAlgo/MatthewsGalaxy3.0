package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestToggleLike_NotAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/posts/test-slug/like", nil)
	c.Params = gin.Params{{Key: "slug", Value: "test-slug"}}
	// Not setting userID in context

	ToggleLike(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// Note: GetLikeStatus tests require database connection
// These would be better suited for integration tests with a test database
