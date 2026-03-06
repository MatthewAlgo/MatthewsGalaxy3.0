package main

import (
	"log"
	"net/http"
	"os"

	"github.com/matthewsgalaxy/email-service/internal/scheduler"
)

func main() {
	log.Println("Starting Matthew's Galaxy Email Service")

	// Create scheduler
	sched, err := scheduler.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	defer sched.Close()

	// Start scheduler in background
	go sched.Start()

	// Simple health check server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"status": "ok", "service": "matthewsgalaxy-email"}`)); err != nil {
			log.Printf("Failed to write health response: %v", err)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Email service health endpoint listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start health server: %v", err)
	}
}
