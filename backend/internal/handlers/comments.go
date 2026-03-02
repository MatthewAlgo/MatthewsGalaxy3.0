package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/database"
	"github.com/matthewsgalaxy/backend/internal/models"
)

// GetComments returns comments for a post
func GetComments(c *gin.Context) {
	postSlug := c.Param("slug")

	// First get the post ID from slug
	var postID uuid.UUID
	err := database.DB.Get(&postID, "SELECT id FROM posts WHERE slug = $1", postSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var comments []models.CommentWithUser
	query := `
		SELECT c.*, u.name as user_name, u.avatar_url as user_avatar_url
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	if err := database.DB.Select(&comments, query, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment adds a comment to a post (authenticated users only)
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	postSlug := c.Param("slug")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get post ID from slug
	var postID uuid.UUID
	err := database.DB.Get(&postID, "SELECT id FROM posts WHERE slug = $1 AND published = true", postSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var comment models.Comment
	err = database.DB.QueryRowx(
		`INSERT INTO comments (post_id, user_id, content)
		 VALUES ($1, $2, $3)
		 RETURNING *`,
		postID, userID, req.Content,
	).StructScan(&comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Get user info for response
	var user models.User
	if err := database.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
		// Comment was created but we can't enrich with user data — return comment alone
		c.JSON(http.StatusCreated, comment)
		return
	}

	response := models.CommentWithUser{
		Comment:       comment,
		UserName:      user.Name,
		UserAvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusCreated, response)
}

// DeleteComment removes a comment (author or admin only)
func DeleteComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("userRole")

	idStr := c.Param("commentId")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Check if user owns the comment or is admin
	var comment models.Comment
	err = database.DB.Get(&comment, "SELECT * FROM comments WHERE id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID.(uuid.UUID) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
