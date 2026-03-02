package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matthewsgalaxy/backend/internal/database"
	"github.com/matthewsgalaxy/backend/internal/handlers"
	"github.com/matthewsgalaxy/backend/internal/middleware"
)

func main() {
	// Validate critical environment variables (fail-fast)
	middleware.InitJWTSecret()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://matthewsgalaxy.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "matthewsgalaxy-backend"})
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Posts (public read, admin write)
		posts := api.Group("/posts")
		{
			posts.GET("", handlers.GetPosts)
			posts.GET("/:slug", handlers.GetPostBySlug)
			posts.GET("/:slug/comments", handlers.GetComments)
			posts.GET("/:slug/likes", middleware.OptionalAuthMiddleware(), handlers.GetLikeStatus)
		}

		// Subscription (public)
		api.POST("/subscribe", middleware.OptionalAuthMiddleware(), handlers.Subscribe)
		api.GET("/unsubscribe", handlers.Unsubscribe)

		// Protected routes (requires authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/me", handlers.GetCurrentUser)
			protected.PATCH("/me", handlers.UpdateProfile)

			// Post interactions
			protected.POST("/posts/:slug/comments", handlers.CreateComment)
			protected.DELETE("/posts/:slug/comments/:commentId", handlers.DeleteComment)
			protected.POST("/posts/:slug/like", handlers.ToggleLike)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// Dashboard
			admin.GET("/stats", handlers.GetDashboardStats)

			// Posts management
			admin.GET("/posts", handlers.GetAllPosts)
			admin.GET("/posts/:id", handlers.GetPostByID)
			admin.POST("/posts", handlers.CreatePost)
			admin.PATCH("/posts/:id", handlers.UpdatePost)
			admin.DELETE("/posts/:id", handlers.DeletePost)

			// Users management
			admin.GET("/users", handlers.GetAllUsers)
			admin.DELETE("/users/:id", handlers.DeleteUser)
			admin.PATCH("/users/:id/role", handlers.UpdateUserRole)

			// Subscribers
			admin.GET("/subscribers", handlers.GetSubscribers)

			// Email logs
			admin.GET("/email-logs", handlers.GetEmailLogs)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Matthew's Galaxy backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
