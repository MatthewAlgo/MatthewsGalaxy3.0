package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/database"
	"github.com/matthewsgalaxy/backend/internal/models"
)

// GetDashboardStats returns statistics for admin dashboard
func GetDashboardStats(c *gin.Context) {
	var stats models.DashboardStats

	database.DB.Get(&stats.TotalUsers, "SELECT COUNT(*) FROM users")
	database.DB.Get(&stats.TotalPosts, "SELECT COUNT(*) FROM posts")
	database.DB.Get(&stats.TotalComments, "SELECT COUNT(*) FROM comments")
	database.DB.Get(&stats.TotalLikes, "SELECT COUNT(*) FROM likes")
	database.DB.Get(&stats.TotalSubscribers, "SELECT COUNT(*) FROM subscriptions WHERE active = true")

	c.JSON(http.StatusOK, stats)
}

// GetAllUsers returns all users for admin
func GetAllUsers(c *gin.Context) {
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		page = max(1, parseInt(p))
	}
	if l := c.Query("limit"); l != "" {
		limit = min(100, max(1, parseInt(l)))
	}

	offset := (page - 1) * limit

	var total int
	database.DB.Get(&total, "SELECT COUNT(*) FROM users")

	var users []models.User
	err := database.DB.Select(&users,
		"SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Convert to response format
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// DeleteUser deletes a user (admin only, cannot delete self or other admins)
func DeleteUser(c *gin.Context) {
	adminID, _ := c.Get("userID")

	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent self-deletion
	if userID == adminID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own account"})
		return
	}

	// Check if target user is admin
	var targetRole string
	err = database.DB.Get(&targetRole, "SELECT role FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if targetRole == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete admin users"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserRole updates a user's role (admin only)
func UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("UPDATE users SET role = $1 WHERE id = $2", req.Role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}
