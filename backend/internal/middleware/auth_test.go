package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Helper to create a valid JWT token for testing
func createTestToken(userID uuid.UUID, email, role string, expired bool) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	if expired {
		expirationTime = time.Now().Add(-24 * time.Hour)
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(getJWTSecret())
	return tokenString
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "user", false)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	var capturedUserID uuid.UUID
	var capturedEmail, capturedRole string

	handler := AuthMiddleware()
	handler(c)

	// Check if context was set
	if id, exists := c.Get("userID"); exists {
		capturedUserID = id.(uuid.UUID)
	}
	if email, exists := c.Get("userEmail"); exists {
		capturedEmail = email.(string)
	}
	if role, exists := c.Get("userRole"); exists {
		capturedRole = role.(string)
	}

	if capturedUserID != userID {
		t.Errorf("UserID mismatch: got %v, want %v", capturedUserID, userID)
	}
	if capturedEmail != "test@example.com" {
		t.Errorf("Email mismatch: got %v, want %v", capturedEmail, "test@example.com")
	}
	if capturedRole != "user" {
		t.Errorf("Role mismatch: got %v, want %v", capturedRole, "user")
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat token123")

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "user", true)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOptionalAuthMiddleware_WithValidToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "admin", false)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := OptionalAuthMiddleware()
	handler(c)

	// Should set context values
	if _, exists := c.Get("userID"); !exists {
		t.Error("Expected userID to be set in context")
	}
	// Should not abort
	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}
}

func TestOptionalAuthMiddleware_WithoutToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler := OptionalAuthMiddleware()
	handler(c)

	// Should not set context values
	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set in context")
	}
	// Should not abort
	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}
}

func TestOptionalAuthMiddleware_WithInvalidToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token")

	handler := OptionalAuthMiddleware()
	handler(c)

	// Should not set context values
	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set in context with invalid token")
	}
	// Should not abort
	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}
}

func TestAdminMiddleware_AdminRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("userRole", "admin")

	handler := AdminMiddleware()
	handler(c)

	if c.IsAborted() {
		t.Error("Admin request should not be aborted")
	}
}

func TestAdminMiddleware_NonAdminRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("userRole", "user")

	handler := AdminMiddleware()
	handler(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestAdminMiddleware_MissingRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	handler := AdminMiddleware()
	handler(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetUserID_Exists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	expectedID := uuid.New()
	c.Set("userID", expectedID)

	id, exists := GetUserID(c)

	if !exists {
		t.Error("Expected exists to be true")
	}
	if id != expectedID {
		t.Errorf("ID mismatch: got %v, want %v", id, expectedID)
	}
}

func TestGetUserID_NotExists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	id, exists := GetUserID(c)

	if exists {
		t.Error("Expected exists to be false")
	}
	if id != uuid.Nil {
		t.Errorf("Expected nil UUID, got %v", id)
	}
}

func TestGetUserRole_Exists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userRole", "admin")

	role := GetUserRole(c)

	if role != "admin" {
		t.Errorf("Role mismatch: got %v, want %v", role, "admin")
	}
}

func TestGetUserRole_NotExists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	role := GetUserRole(c)

	if role != "" {
		t.Errorf("Expected empty role, got %v", role)
	}
}

func TestGetJWTSecret_Default(t *testing.T) {
	// Clear environment variable
	os.Unsetenv("JWT_SECRET")

	secret := getJWTSecret()

	if string(secret) != "your-secret-key" {
		t.Errorf("Expected default secret, got %s", string(secret))
	}
}

func TestGetJWTSecret_FromEnv(t *testing.T) {
	os.Setenv("JWT_SECRET", "custom_secret_123")
	defer os.Unsetenv("JWT_SECRET")

	secret := getJWTSecret()

	if string(secret) != "custom_secret_123" {
		t.Errorf("Expected custom secret, got %s", string(secret))
	}
}
