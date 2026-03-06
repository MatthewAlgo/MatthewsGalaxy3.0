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

	if err := database.DB.Get(&stats.TotalUsers, "SELECT COUNT(*) FROM users"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
		return
	}
	if err := database.DB.Get(&stats.TotalPosts, "SELECT COUNT(*) FROM posts"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post stats"})
		return
	}
	if err := database.DB.Get(&stats.TotalComments, "SELECT COUNT(*) FROM comments"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment stats"})
		return
	}
	if err := database.DB.Get(&stats.TotalLikes, "SELECT COUNT(*) FROM likes"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch like stats"})
		return
	}
	if err := database.DB.Get(&stats.TotalSubscribers, "SELECT COUNT(*) FROM subscriptions WHERE active = true"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriber stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAllUsers returns all users for admin
func GetAllUsers(c *gin.Context) {
	page, limit, offset := parsePagination(c, 20)

	var total int
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM users"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

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

	c.JSON(http.StatusOK, buildPaginatedResponse(responses, page, limit, total))
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

// GetEmailLogs returns paginated email delivery logs (admin only)
func GetEmailLogs(c *gin.Context) {
	page, limit, offset := parsePagination(c, 50)

	var total int
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM email_logs"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count email logs"})
		return
	}

	var logs []models.EmailLog
	err := database.DB.Select(&logs,
		"SELECT * FROM email_logs ORDER BY sent_at DESC LIMIT $1 OFFSET $2",
		limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email logs"})
		return
	}

	c.JSON(http.StatusOK, buildPaginatedResponse(logs, page, limit, total))
}
