package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/database"
)

// ToggleLike adds or removes a like from a post
func ToggleLike(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	postSlug := c.Param("slug")

	// Get post ID from slug
	var postID uuid.UUID
	err := database.DB.Get(&postID, "SELECT id FROM posts WHERE slug = $1 AND published = true", postSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if like exists
	var exists_like bool
	err = database.DB.Get(&exists_like,
		"SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2)",
		postID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists_like {
		// Remove like
		_, err = database.DB.Exec("DELETE FROM likes WHERE post_id = $1 AND user_id = $2", postID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"liked": false, "message": "Like removed"})
	} else {
		// Add like
		_, err = database.DB.Exec("INSERT INTO likes (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"liked": true, "message": "Post liked"})
	}
}

// GetLikeStatus checks if current user has liked a post
func GetLikeStatus(c *gin.Context) {
	postSlug := c.Param("slug")

	// Get post ID from slug
	var postID uuid.UUID
	err := database.DB.Get(&postID, "SELECT id FROM posts WHERE slug = $1", postSlug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Get like count
	var likeCount int
	if err := database.DB.Get(&likeCount, "SELECT COUNT(*) FROM likes WHERE post_id = $1", postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch like count"})
		return
	}

	response := gin.H{
		"count": likeCount,
		"liked": false,
	}

	// Check if user is authenticated and has liked
	userID, exists := c.Get("userID")
	if exists {
		var liked bool
		if err := database.DB.Get(&liked,
			"SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2)",
			postID, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
			return
		}
		response["liked"] = liked
	}

	c.JSON(http.StatusOK, response)
}
