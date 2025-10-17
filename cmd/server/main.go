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

	// Run admin features migration
	adminMigrationPath := "migrations/002_add_admin_features.sql"
	if err := db.Migrate(adminMigrationPath); err != nil {
		log.Printf("Warning: Admin migration failed: %v", err)
	}

	// Initialize handlers
	h := handlers.New(db)
	adminH := handlers.NewAdminHandler(db)

	// Create router
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", h.HealthHandler)
	mux.HandleFunc("/api/register", h.RegisterHandler)
	mux.HandleFunc("/api/lookup", h.LookupHandler)
	mux.HandleFunc("/api/search", h.SearchHandler)

	// Protected API routes (require authentication)
	authMiddleware := handlers.AuthMiddleware(db)
	mux.Handle("/api/upload", authMiddleware(http.HandlerFunc(h.UploadHandler)))
	mux.Handle("/api/vote", authMiddleware(http.HandlerFunc(h.VoteHandler)))
	mux.Handle("/api/vote/delete", authMiddleware(http.HandlerFunc(h.DeleteVoteHandler)))

	// Admin API routes (require authentication + admin role)
	adminMiddleware := func(handler http.HandlerFunc) http.Handler {
		return authMiddleware(handlers.AdminMiddleware(http.HandlerFunc(handler)))
	}
	mux.Handle("/api/admin/submissions", adminMiddleware(adminH.ListSubmissionsHandler))
	mux.Handle("/api/admin/submissions/get", adminMiddleware(adminH.GetSubmissionHandler))
	mux.Handle("/api/admin/submissions/delete", adminMiddleware(adminH.DeleteSubmissionHandler))
	mux.Handle("/api/admin/users", adminMiddleware(adminH.ListUsersHandler))
	mux.Handle("/api/admin/users/get", adminMiddleware(adminH.GetUserHandler))
	mux.Handle("/api/admin/users/role", adminMiddleware(adminH.ChangeUserRoleHandler))
	mux.Handle("/api/admin/users/status", adminMiddleware(adminH.ChangeUserStatusHandler))
	mux.Handle("/api/admin/stats", adminMiddleware(adminH.GetStatsHandler))

	// Serve static frontend files
	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "frontend/public"
	}

	// Check if frontend directory exists
	if _, err := os.Stat(frontendPath); err == nil {
		log.Printf("Serving frontend from: %s", frontendPath)
		fs := http.FileServer(http.Dir(frontendPath))
		mux.Handle("/", fs)
	} else {
		log.Printf("Frontend not found at %s, serving API only", frontendPath)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"MKV Mender API","endpoints":["/api/health","/api/register","/api/lookup","/api/search","/api/upload","/api/vote"]}`))
		})
	}

	// Apply global middleware
	handler := handlers.LoggingMiddleware(handlers.CORSMiddleware(mux))

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
