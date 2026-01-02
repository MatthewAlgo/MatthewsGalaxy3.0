package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	Role         string    `db:"role" json:"role"`
	AvatarURL    *string   `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio          *string   `db:"bio" json:"bio,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// UserResponse is what we return to clients (no sensitive data)
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Bio       *string   `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
	}
}

// Post represents a blog post
type Post struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Slug       string    `db:"slug" json:"slug"`
	Content    string    `db:"content" json:"content"`
	Excerpt    *string   `db:"excerpt" json:"excerpt,omitempty"`
	CoverImage *string   `db:"cover_image" json:"cover_image,omitempty"`
	AuthorID   uuid.UUID `db:"author_id" json:"author_id"`
	Published  bool      `db:"published" json:"published"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// PostWithAuthor includes author information
type PostWithAuthor struct {
	Post
	AuthorName      string  `db:"author_name" json:"author_name"`
	AuthorAvatarURL *string `db:"author_avatar_url" json:"author_avatar_url,omitempty"`
	LikeCount       int     `db:"like_count" json:"like_count"`
	CommentCount    int     `db:"comment_count" json:"comment_count"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PostID    uuid.UUID `db:"post_id" json:"post_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CommentWithUser includes user information
type CommentWithUser struct {
	Comment
	UserName      string  `db:"user_name" json:"user_name"`
	UserAvatarURL *string `db:"user_avatar_url" json:"user_avatar_url,omitempty"`
}

// Like represents a like on a post
type Like struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PostID    uuid.UUID `db:"post_id" json:"post_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Subscription represents an email subscription
type Subscription struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         *uuid.UUID `db:"user_id" json:"user_id,omitempty"`
	Email          string     `db:"email" json:"email"`
	Active         bool       `db:"active" json:"active"`
	SubscribedAt   time.Time  `db:"subscribed_at" json:"subscribed_at"`
	UnsubscribedAt *time.Time `db:"unsubscribed_at" json:"unsubscribed_at,omitempty"`
}

// EmailLog tracks sent emails
type EmailLog struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	PostID          *uuid.UUID `db:"post_id" json:"post_id,omitempty"`
	SubscriberEmail string     `db:"subscriber_email" json:"subscriber_email"`
	SentAt          time.Time  `db:"sent_at" json:"sent_at"`
	Status          string     `db:"status" json:"status"`
}

// Request/Response DTOs

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CreatePostRequest struct {
	Title      string  `json:"title" binding:"required,min=3"`
	Content    string  `json:"content" binding:"required,min=10"`
	Excerpt    *string `json:"excerpt"`
	CoverImage *string `json:"cover_image"`
	Published  bool    `json:"published"`
}

type UpdatePostRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Excerpt    *string `json:"excerpt"`
	CoverImage *string `json:"cover_image"`
	Published  *bool   `json:"published"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PaginationParams struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// Stats for admin dashboard
type DashboardStats struct {
	TotalUsers       int `json:"total_users"`
	TotalPosts       int `json:"total_posts"`
	TotalComments    int `json:"total_comments"`
	TotalLikes       int `json:"total_likes"`
	TotalSubscribers int `json:"total_subscribers"`
}
