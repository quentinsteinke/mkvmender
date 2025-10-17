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
		INSERT INTO users (username, api_key, role, is_active)
		VALUES (?, ?, 'user', 1)
		RETURNING id, username, api_key, role, is_active, created_at, updated_at
	`

	var user models.User
	err = db.conn.QueryRow(query, username, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.Role,
		&user.IsActive,
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
		SELECT id, username, api_key, role, is_active, created_at, updated_at
		FROM users
		WHERE api_key = ?
	`

	var user models.User
	err := db.conn.QueryRow(query, apiKey).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.Role,
		&user.IsActive,
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
		SELECT id, username, api_key, role, is_active, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	var user models.User
	err := db.conn.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.Role,
		&user.IsActive,
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

// GetUserByID retrieves a user by their ID
func (db *DB) GetUserByID(userID int64) (*models.User, error) {
	query := `
		SELECT id, username, api_key, role, is_active, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := db.conn.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.APIKey,
		&user.Role,
		&user.IsActive,
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

// UpdateUserRole updates a user's role
func (db *DB) UpdateUserRole(userID int64, role models.UserRole) error {
	query := `
		UPDATE users
		SET role = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := db.conn.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateUserStatus updates a user's active status
func (db *DB) UpdateUserStatus(userID int64, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := db.conn.Exec(query, isActive, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// generateAPIKey generates a random 32-byte API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
