package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/models"
)

func TestCreateComment_NotAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/posts/test-slug/comments", nil)
	c.Params = gin.Params{{Key: "slug", Value: "test-slug"}}
	// Not setting userID in context

	CreateComment(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreateComment_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/posts/test-slug/comments", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "slug", Value: "test-slug"}}
	c.Set("userID", uuid.New())

	CreateComment(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateComment_EmptyContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := models.CreateCommentRequest{
		Content: "", // Empty content should fail validation
	}
	body, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/api/v1/posts/test-slug/comments", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "slug", Value: "test-slug"}}
	c.Set("userID", uuid.New())

	CreateComment(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteComment_InvalidCommentID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/posts/test-slug/comments/invalid-uuid", nil)
	c.Params = gin.Params{
		{Key: "slug", Value: "test-slug"},
		{Key: "commentId", Value: "not-a-valid-uuid"},
	}
	c.Set("userID", uuid.New())
	c.Set("userRole", "user")

	DeleteComment(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
