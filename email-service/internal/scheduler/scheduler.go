package scheduler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/matthewsgalaxy/email-service/internal/sender"
	"github.com/matthewsgalaxy/email-service/internal/templates"
)

type Post struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Slug      string    `db:"slug"`
	Excerpt   *string   `db:"excerpt"`
	CreatedAt time.Time `db:"created_at"`
}

type Subscription struct {
	ID    uuid.UUID `db:"id"`
	Email string    `db:"email"`
}

type Scheduler struct {
	db           *sqlx.DB
	emailSender  *sender.EmailSender
	scanInterval time.Duration
	lastScanTime time.Time
	frontendURL  string
	backendURL   string
}

// NewScheduler creates a new post scanner and email scheduler
func NewScheduler() (*Scheduler, error) {
	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/matthewsgalaxy?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Parse scan interval
	intervalSeconds, _ := strconv.Atoi(os.Getenv("SCAN_INTERVAL"))
	if intervalSeconds == 0 {
		intervalSeconds = 3600 // Default 1 hour
	}

	return &Scheduler{
		db:           db,
		emailSender:  sender.NewEmailSender(),
		scanInterval: time.Duration(intervalSeconds) * time.Second,
		lastScanTime: time.Now().Add(-time.Duration(intervalSeconds) * time.Second),
		frontendURL:  getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		backendURL:   getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
	}, nil
}

// Start begins the email notification scheduler
func (s *Scheduler) Start() {
	log.Printf("Starting email scheduler with %v scan interval", s.scanInterval)

	// Run immediately on start
	s.scanAndNotify()

	// Then run on interval
	ticker := time.NewTicker(s.scanInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.scanAndNotify()
	}
}

// scanAndNotify checks for new posts and sends notifications
func (s *Scheduler) scanAndNotify() {
	log.Printf("Scanning for new posts since %v", s.lastScanTime)

	// Get new published posts since last scan
	var newPosts []Post
	err := s.db.Select(&newPosts, `
		SELECT id, title, slug, excerpt, created_at 
		FROM posts 
		WHERE published = true AND created_at > $1
		ORDER BY created_at ASC
	`, s.lastScanTime)

	if err != nil {
		log.Printf("Error fetching new posts: %v", err)
		return
	}

	if len(newPosts) == 0 {
		log.Println("No new posts found")
		s.lastScanTime = time.Now()
		return
	}

	log.Printf("Found %d new posts", len(newPosts))

	// Get active subscribers
	var subscribers []Subscription
	err = s.db.Select(&subscribers, "SELECT id, email FROM subscriptions WHERE active = true")
	if err != nil {
		log.Printf("Error fetching subscribers: %v", err)
		return
	}

	if len(subscribers) == 0 {
		log.Println("No active subscribers")
		s.lastScanTime = time.Now()
		return
	}

	log.Printf("Sending notifications to %d subscribers", len(subscribers))

	// Send notification for each new post
	for _, post := range newPosts {
		s.notifySubscribers(post, subscribers)
	}

	s.lastScanTime = time.Now()
}

// notifySubscribers sends email notifications for a post
func (s *Scheduler) notifySubscribers(post Post, subscribers []Subscription) {
	postURL := fmt.Sprintf("%s/blog/%s", s.frontendURL, post.Slug)

	excerpt := "Check out the latest post on Matthew's Galaxy!"
	if post.Excerpt != nil && *post.Excerpt != "" {
		excerpt = *post.Excerpt
	}

	for _, sub := range subscribers {
		unsubscribeURL := fmt.Sprintf("%s/api/v1/unsubscribe?email=%s", s.backendURL, sub.Email)

		htmlBody, err := templates.NewPostEmail(post.Title, excerpt, postURL, unsubscribeURL)
		if err != nil {
			log.Printf("Error generating email template: %v", err)
			continue
		}

		subject := fmt.Sprintf("✨ New Post: %s - Matthew's Galaxy", post.Title)

		if err := s.emailSender.SendEmail(sub.Email, subject, htmlBody); err != nil {
			log.Printf("Failed to send email to %s: %v", sub.Email, err)
			s.logEmailSent(post.ID, sub.Email, "failed")
		} else {
			s.logEmailSent(post.ID, sub.Email, "sent")
		}

		// Small delay between emails to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}
}

// logEmailSent records the email in the database
func (s *Scheduler) logEmailSent(postID uuid.UUID, email, status string) {
	_, err := s.db.Exec(`
		INSERT INTO email_logs (post_id, subscriber_email, status)
		VALUES ($1, $2, $3)
	`, postID, email, status)
	if err != nil {
		log.Printf("Failed to log email: %v", err)
	}
}

// Close closes the database connection
func (s *Scheduler) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
