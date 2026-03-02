package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParsePagination_Defaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	page, limit, offset := parsePagination(c, 10)

	if page != 1 {
		t.Errorf("Expected page 1, got %d", page)
	}
	if limit != 10 {
		t.Errorf("Expected limit 10, got %d", limit)
	}
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}
}

func TestParsePagination_CustomValues(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=3&limit=25", nil)

	page, limit, offset := parsePagination(c, 10)

	if page != 3 {
		t.Errorf("Expected page 3, got %d", page)
	}
	if limit != 25 {
		t.Errorf("Expected limit 25, got %d", limit)
	}
	if offset != 50 {
		t.Errorf("Expected offset 50, got %d", offset)
	}
}

func TestParsePagination_NegativePage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=0", nil)

	page, _, _ := parsePagination(c, 10)

	if page != 1 {
		t.Errorf("Expected page clamped to 1, got %d", page)
	}
}

func TestParsePagination_LimitClamped(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?limit=500", nil)

	_, limit, _ := parsePagination(c, 10)

	if limit != 100 {
		t.Errorf("Expected limit clamped to 100, got %d", limit)
	}
}

func TestParsePagination_ZeroLimit(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?limit=0", nil)

	_, limit, _ := parsePagination(c, 10)

	if limit != 1 {
		t.Errorf("Expected limit clamped to 1, got %d", limit)
	}
}

func TestParsePagination_DifferentDefaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	_, limit, _ := parsePagination(c, 20)

	if limit != 20 {
		t.Errorf("Expected default limit 20, got %d", limit)
	}
}

func TestBuildPaginatedResponse(t *testing.T) {
	data := []string{"a", "b", "c"}
	resp := buildPaginatedResponse(data, 2, 10, 25)

	if resp.Page != 2 {
		t.Errorf("Expected page 2, got %d", resp.Page)
	}
	if resp.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", resp.Limit)
	}
	if resp.Total != 25 {
		t.Errorf("Expected total 25, got %d", resp.Total)
	}
	if resp.TotalPages != 3 {
		t.Errorf("Expected totalPages 3, got %d", resp.TotalPages)
	}
}

func TestBuildPaginatedResponse_ExactDivision(t *testing.T) {
	resp := buildPaginatedResponse(nil, 1, 10, 30)

	if resp.TotalPages != 3 {
		t.Errorf("Expected totalPages 3 for 30/10, got %d", resp.TotalPages)
	}
}

func TestBuildPaginatedResponse_SinglePage(t *testing.T) {
	resp := buildPaginatedResponse(nil, 1, 10, 5)

	if resp.TotalPages != 1 {
		t.Errorf("Expected totalPages 1 for 5/10, got %d", resp.TotalPages)
	}
}

func TestBuildPaginatedResponse_ZeroResults(t *testing.T) {
	resp := buildPaginatedResponse(nil, 1, 10, 0)

	if resp.TotalPages != 0 {
		t.Errorf("Expected totalPages 0, got %d", resp.TotalPages)
	}
}
