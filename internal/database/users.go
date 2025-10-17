package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// CreateUser creates a new user with a generated API key
func (db *DB) CreateUser(username string) (*models.User, error) {
	// Generate random API key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	query := `
		INSERT INTO users (username, api_key)
		VALUES (?, ?)
		RETURNING id, username, api_key, created_at, updated_at
	`

	var user models.User
	err = db.conn.QueryRow(query, username, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByAPIKey retrieves a user by their API key
func (db *DB) GetUserByAPIKey(apiKey string) (*models.User, error) {
	query := `
		SELECT id, username, api_key, created_at, updated_at
		FROM users
		WHERE api_key = ?
	`

	var user models.User
	err := db.conn.QueryRow(query, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, api_key, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	var user models.User
	err := db.conn.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// generateAPIKey generates a random 32-byte API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
