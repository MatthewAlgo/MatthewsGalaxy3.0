package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserToResponse(t *testing.T) {
	avatarURL := "https://example.com/avatar.jpg"
	bio := "Software Engineer"
	now := time.Now()

	user := User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "secret_hash_should_not_appear",
		Name:         "Test User",
		Role:         "user",
		AvatarURL:    &avatarURL,
		Bio:          &bio,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	response := user.ToResponse()

	// Verify all fields are correctly mapped
	if response.ID != user.ID {
		t.Errorf("ID mismatch: got %v, want %v", response.ID, user.ID)
	}
	if response.Email != user.Email {
		t.Errorf("Email mismatch: got %v, want %v", response.Email, user.Email)
	}
	if response.Name != user.Name {
		t.Errorf("Name mismatch: got %v, want %v", response.Name, user.Name)
	}
	if response.Role != user.Role {
		t.Errorf("Role mismatch: got %v, want %v", response.Role, user.Role)
	}
	if response.AvatarURL == nil || *response.AvatarURL != avatarURL {
		t.Errorf("AvatarURL mismatch: got %v, want %v", response.AvatarURL, avatarURL)
	}
	if response.Bio == nil || *response.Bio != bio {
		t.Errorf("Bio mismatch: got %v, want %v", response.Bio, bio)
	}
	if response.CreatedAt != user.CreatedAt {
		t.Errorf("CreatedAt mismatch: got %v, want %v", response.CreatedAt, user.CreatedAt)
	}
}

func TestUserToResponseWithNilOptionalFields(t *testing.T) {
	now := time.Now()

	user := User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "secret_hash",
		Name:         "Test User",
		Role:         "admin",
		AvatarURL:    nil,
		Bio:          nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	response := user.ToResponse()

	if response.AvatarURL != nil {
		t.Errorf("AvatarURL should be nil, got %v", response.AvatarURL)
	}
	if response.Bio != nil {
		t.Errorf("Bio should be nil, got %v", response.Bio)
	}
}

func TestPaginationParamsDefaults(t *testing.T) {
	params := PaginationParams{
		Page:  1,
		Limit: 10,
	}

	if params.Page != 1 {
		t.Errorf("Page should be 1, got %d", params.Page)
	}
	if params.Limit != 10 {
		t.Errorf("Limit should be 10, got %d", params.Limit)
	}
}

func TestPaginatedResponseStructure(t *testing.T) {
	data := []string{"item1", "item2"}
	response := PaginatedResponse{
		Data:       data,
		Page:       1,
		Limit:      10,
		Total:      100,
		TotalPages: 10,
	}

	if response.Page != 1 {
		t.Errorf("Page should be 1, got %d", response.Page)
	}
	if response.Total != 100 {
		t.Errorf("Total should be 100, got %d", response.Total)
	}
	if response.TotalPages != 10 {
		t.Errorf("TotalPages should be 10, got %d", response.TotalPages)
	}
}

func TestDashboardStatsStructure(t *testing.T) {
	stats := DashboardStats{
		TotalUsers:       50,
		TotalPosts:       25,
		TotalComments:    100,
		TotalLikes:       200,
		TotalSubscribers: 75,
	}

	if stats.TotalUsers != 50 {
		t.Errorf("TotalUsers should be 50, got %d", stats.TotalUsers)
	}
	if stats.TotalPosts != 25 {
		t.Errorf("TotalPosts should be 25, got %d", stats.TotalPosts)
	}
}
