package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matthewsgalaxy/backend/internal/models"
)

func TestSubscribe_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/subscribe", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	Subscribe(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSubscribe_InvalidEmail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := models.SubscribeRequest{
		Email: "not-a-valid-email",
	}
	body, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/api/v1/subscribe", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Subscribe(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSubscribe_MissingEmail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := models.SubscribeRequest{
		Email: "",
	}
	body, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/api/v1/subscribe", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Subscribe(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUnsubscribe_MissingEmail(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/unsubscribe", nil)
	// No email query param

	Unsubscribe(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Note: TestUnsubscribe_WithEmailParam requires database connection
// This would be better suited for integration tests with a test database
