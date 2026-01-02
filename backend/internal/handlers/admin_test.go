package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestDeleteUser_InvalidUUID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/users/invalid-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "not-a-valid-uuid"}}
	c.Set("userID", uuid.New())

	DeleteUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteUser_SelfDeletion(t *testing.T) {
	adminID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/v1/admin/users/"+adminID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: adminID.String()}}
	c.Set("userID", adminID)

	DeleteUser(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestUpdateUserRole_InvalidUUID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/users/invalid-uuid/role", nil)
	c.Params = gin.Params{{Key: "id", Value: "not-a-valid-uuid"}}

	UpdateUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserRole_InvalidJSON(t *testing.T) {
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/users/"+userID.String()+"/role", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	UpdateUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserRole_InvalidRole(t *testing.T) {
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := []byte(`{"role": "superadmin"}`) // Invalid role
	c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/users/"+userID.String()+"/role", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	UpdateUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserRole_MissingRole(t *testing.T) {
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := []byte(`{}`) // Missing role
	c.Request = httptest.NewRequest("PATCH", "/api/v1/admin/users/"+userID.String()+"/role", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	UpdateUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
