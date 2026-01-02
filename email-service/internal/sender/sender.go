package sender

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewEmailSender creates a new email sender
func NewEmailSender() *EmailSender {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}

	return &EmailSender{
		host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		port:     port,
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     getEnvOrDefault("FROM_EMAIL", "noreply@matthewsgalaxy.com"),
	}
}

// SendEmail sends an HTML email
func (s *EmailSender) SendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// SendBulkEmails sends emails to multiple recipients
func (s *EmailSender) SendBulkEmails(recipients []string, subject string, htmlBodyFunc func(email string) string) []error {
	var errors []error

	for _, to := range recipients {
		htmlBody := htmlBodyFunc(to)
		if err := s.SendEmail(to, subject, htmlBody); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
