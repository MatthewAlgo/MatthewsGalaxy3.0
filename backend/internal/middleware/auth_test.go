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
	// Set a test secret for all tests
	os.Setenv("JWT_SECRET", "test-secret-key-that-is-long-enough-for-testing")
	InitJWTSecret()
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
	tokenString, _ := token.SignedString(GetJWTSecret())
	return tokenString
}

// ==========================================
// AuthMiddleware Tests
// ==========================================

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "user", false)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := AuthMiddleware()
	handler(c)

	// Check if context was set correctly
	if id, exists := c.Get("userID"); !exists {
		t.Error("Expected userID to be set in context")
	} else if id.(uuid.UUID) != userID {
		t.Errorf("UserID mismatch: got %v, want %v", id, userID)
	}

	if email, exists := c.Get("userEmail"); !exists || email != "test@example.com" {
		t.Errorf("Email mismatch: got %v", email)
	}

	if role, exists := c.Get("userRole"); !exists || role != "user" {
		t.Errorf("Role mismatch: got %v", role)
	}

	if c.IsAborted() {
		t.Error("Request should not be aborted for valid token")
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
	if !c.IsAborted() {
		t.Error("Request should be aborted for missing header")
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

func TestAuthMiddleware_TokenSignedWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("wrong-secret-key-that-is-different"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d for wrong-secret token, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_EmptyBearerToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	handler := AuthMiddleware()
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// ==========================================
// OptionalAuthMiddleware Tests
// ==========================================

func TestOptionalAuthMiddleware_WithValidToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "admin", false)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := OptionalAuthMiddleware()
	handler(c)

	if _, exists := c.Get("userID"); !exists {
		t.Error("Expected userID to be set in context")
	}
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

	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set in context")
	}
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

	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set in context with invalid token")
	}
	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}
}

func TestOptionalAuthMiddleware_WithExpiredToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "user", true)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	handler := OptionalAuthMiddleware()
	handler(c)

	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set for expired token")
	}
	if c.IsAborted() {
		t.Error("Request should not be aborted for expired token in optional middleware")
	}
}

func TestOptionalAuthMiddleware_WithInvalidFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Basic sometoken")

	handler := OptionalAuthMiddleware()
	handler(c)

	if _, exists := c.Get("userID"); exists {
		t.Error("Expected userID to NOT be set for non-Bearer format")
	}
	if c.IsAborted() {
		t.Error("Request should not be aborted")
	}
}

// ==========================================
// AdminMiddleware Tests
// ==========================================

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

// ==========================================
// Helper Function Tests
// ==========================================

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

// ==========================================
// GenerateToken + ParseToken Tests
// ==========================================

func TestGenerateToken_CreatesValidToken(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	role := "admin"

	token, err := GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Verify roundtrip through ParseToken
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email mismatch: got %v, want %v", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("Role mismatch: got %v, want %v", claims.Role, role)
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	token := createTestToken(userID, "test@example.com", "user", true)

	_, err := ParseToken(token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestParseToken_InvalidSignature(t *testing.T) {
	claims := &Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("wrong-secret"))

	_, err := ParseToken(tokenString)
	if err == nil {
		t.Error("Expected error for token signed with wrong secret")
	}
}

func TestParseToken_MalformedToken(t *testing.T) {
	_, err := ParseToken("not-a-valid-jwt")
	if err == nil {
		t.Error("Expected error for malformed token")
	}
}

func TestParseToken_EmptyToken(t *testing.T) {
	_, err := ParseToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestGetJWTSecret_ReturnsInitializedSecret(t *testing.T) {
	secret := GetJWTSecret()
	if string(secret) != "test-secret-key-that-is-long-enough-for-testing" {
		t.Errorf("Expected test secret, got %s", string(secret))
	}
}
