package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/database"
	"github.com/matthewsgalaxy/backend/internal/models"
)

// Subscribe adds an email to the subscription list
func Subscribe(c *gin.Context) {
	var req models.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already subscribed
	var existing models.Subscription
	err := database.DB.Get(&existing, "SELECT * FROM subscriptions WHERE email = $1", req.Email)

	if err == nil {
		// Already exists
		if existing.Active {
			c.JSON(http.StatusOK, gin.H{"message": "Already subscribed"})
			return
		}
		// Reactivate subscription
		_, err = database.DB.Exec(
			"UPDATE subscriptions SET active = true, unsubscribed_at = NULL WHERE email = $1",
			req.Email,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reactivate subscription"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Subscription reactivated"})
		return
	}

	// Get user ID if authenticated
	var userID *uuid.UUID
	if id, exists := c.Get("userID"); exists {
		uid := id.(uuid.UUID)
		userID = &uid
	}

	// Create new subscription
	_, err = database.DB.Exec(
		"INSERT INTO subscriptions (user_id, email, active) VALUES ($1, $2, true)",
		userID, req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully subscribed"})
}

// Unsubscribe removes an email from the subscription list
func Unsubscribe(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email required"})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE subscriptions SET active = false, unsubscribed_at = NOW() WHERE email = $1",
		email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unsubscribed"})
}

// GetSubscribers returns all active subscribers (admin only)
func GetSubscribers(c *gin.Context) {
	var subscribers []models.Subscription
	err := database.DB.Select(&subscribers,
		"SELECT * FROM subscriptions WHERE active = true ORDER BY subscribed_at DESC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscribers"})
		return
	}

	c.JSON(http.StatusOK, subscribers)
}
