package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quentinsteinke/mkvmender/internal/database"
	"github.com/quentinsteinke/mkvmender/internal/handlers"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	db, err := database.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		migrationPath = "migrations/001_initial_schema.sql"
	}

	if err := db.Migrate(migrationPath); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	// Initialize handlers
	h := handlers.New(db)

	// Create router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/health", h.HealthHandler)
	mux.HandleFunc("/api/register", h.RegisterHandler)
	mux.HandleFunc("/api/lookup", h.LookupHandler)
	mux.HandleFunc("/api/search", h.SearchHandler)

	// Protected routes (require authentication)
	authMiddleware := handlers.AuthMiddleware(db)
	mux.Handle("/api/upload", authMiddleware(http.HandlerFunc(h.UploadHandler)))
	mux.Handle("/api/vote", authMiddleware(http.HandlerFunc(h.VoteHandler)))
	mux.Handle("/api/vote/delete", authMiddleware(http.HandlerFunc(h.DeleteVoteHandler)))

	// Apply global middleware
	handler := handlers.LoggingMiddleware(handlers.CORSMiddleware(mux))

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
