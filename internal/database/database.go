package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// Config holds database configuration
type Config struct {
	URL      string
	AuthToken string
}

// New creates a new database connection
func New(cfg Config) (*DB, error) {
	// For Turso, the DSN format is: libsql://[host]?authToken=[token]
	// or for local: file:[path]
	var dsn string
	if cfg.URL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// If it's a Turso URL, append auth token
	if cfg.AuthToken != "" {
		dsn = fmt.Sprintf("%s?authToken=%s", cfg.URL, cfg.AuthToken)
	} else {
		dsn = cfg.URL
	}

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: db}, nil
}

// NewFromEnv creates a database connection from environment variables
func NewFromEnv() (*DB, error) {
	cfg := Config{
		URL:      os.Getenv("TURSO_DATABASE_URL"),
		AuthToken: os.Getenv("TURSO_AUTH_TOKEN"),
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL environment variable is not set")
	}

	return New(cfg)
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying sql.DB connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Migrate runs database migrations
func (db *DB) Migrate(migrationPath string) error {
	// Read migration file
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	_, err = db.conn.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// Begin starts a new transaction
func (db *DB) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}
