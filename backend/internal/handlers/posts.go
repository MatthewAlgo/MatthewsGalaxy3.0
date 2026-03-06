package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matthewsgalaxy/backend/internal/database"
	"github.com/matthewsgalaxy/backend/internal/models"
)

// parsePagination extracts page, limit, and offset from query parameters.
// This is the single source of truth for pagination logic — do not duplicate.
func parsePagination(c *gin.Context, defaultLimit int) (page, limit, offset int) {
	page = 1
	limit = defaultLimit

	if p := c.Query("page"); p != "" {
		page = max(1, parseInt(p))
	}
	if l := c.Query("limit"); l != "" {
		limit = min(100, max(1, parseInt(l)))
	}

	offset = (page - 1) * limit
	return
}

// buildPaginatedResponse creates a PaginatedResponse with computed total pages.
func buildPaginatedResponse(data interface{}, page, limit, total int) models.PaginatedResponse {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	return models.PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetPosts returns paginated published posts
func GetPosts(c *gin.Context) {
	page, limit, offset := parsePagination(c, 10)

	var total int
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM posts WHERE published = true"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count posts"})
		return
	}

	var posts []models.PostWithAuthor
	query := `
		SELECT p.*, u.name as author_name, u.avatar_url as author_avatar_url,
			   COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id), 0) as like_count,
			   COALESCE((SELECT COUNT(*) FROM comments WHERE post_id = p.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.published = true
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	if err := database.DB.Select(&posts, query, limit, offset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, buildPaginatedResponse(posts, page, limit, total))
}

// GetPostBySlug returns a single post by slug
func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var post models.PostWithAuthor
	query := `
		SELECT p.*, u.name as author_name, u.avatar_url as author_avatar_url,
			   COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id), 0) as like_count,
			   COALESCE((SELECT COUNT(*) FROM comments WHERE post_id = p.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.slug = $1 AND p.published = true
	`

	if err := database.DB.Get(&post, query, slug); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPostByID returns a single post by ID (admin)
func GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.PostWithAuthor
	query := `
		SELECT p.*, u.name as author_name, u.avatar_url as author_avatar_url,
			   COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id), 0) as like_count,
			   COALESCE((SELECT COUNT(*) FROM comments WHERE post_id = p.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`

	if err := database.DB.Get(&post, query, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost creates a new post (admin only)
func CreatePost(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := generateSlug(req.Title)

	// Ensure unique slug
	var count int
	if err := database.DB.Get(&count, "SELECT COUNT(*) FROM posts WHERE slug LIKE $1", slug+"%"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check slug uniqueness"})
		return
	}
	if count > 0 {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	var post models.Post
	err := database.DB.QueryRowx(
		`INSERT INTO posts (title, slug, content, excerpt, cover_image, author_id, published)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING *`,
		req.Title, slug, req.Content, req.Excerpt, req.CoverImage, userID, req.Published,
	).StructScan(&post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost updates an existing post (admin only)
func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	err = database.DB.QueryRowx(
		`UPDATE posts SET
			title = COALESCE($1, title),
			content = COALESCE($2, content),
			excerpt = COALESCE($3, excerpt),
			cover_image = COALESCE($4, cover_image),
			published = COALESCE($5, published),
			updated_at = NOW()
		 WHERE id = $6
		 RETURNING *`,
		req.Title, req.Content, req.Excerpt, req.CoverImage, req.Published, id,
	).StructScan(&post)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a post (admin only)
func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetAllPosts returns all posts for admin
func GetAllPosts(c *gin.Context) {
	page, limit, offset := parsePagination(c, 20)

	var total int
	if err := database.DB.Get(&total, "SELECT COUNT(*) FROM posts"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count posts"})
		return
	}

	var posts []models.PostWithAuthor
	query := `
		SELECT p.*, u.name as author_name, u.avatar_url as author_avatar_url,
			   COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id), 0) as like_count,
			   COALESCE((SELECT COUNT(*) FROM comments WHERE post_id = p.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	if err := database.DB.Select(&posts, query, limit, offset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, buildPaginatedResponse(posts, page, limit, total))
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
