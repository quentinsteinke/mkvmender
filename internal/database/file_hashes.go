package database

import (
	"database/sql"
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// CreateFileHash creates a new file hash entry or returns existing one
func (db *DB) CreateFileHash(hash string, fileSize int64, mediaType models.MediaType) (*models.FileHash, error) {
	// First, check if hash already exists
	existing, err := db.GetFileHashByHash(hash)
	if err == nil {
		return existing, nil
	}

	query := `
		INSERT INTO file_hashes (hash, file_size, media_type)
		VALUES (?, ?, ?)
		RETURNING id, hash, file_size, media_type, created_at
	`

	var fileHash models.FileHash
	err = db.conn.QueryRow(query, hash, fileSize, string(mediaType)).Scan(
		&fileHash.ID,
		&fileHash.Hash,
		&fileHash.FileSize,
		&fileHash.MediaType,
		&fileHash.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create file hash: %w", err)
	}

	return &fileHash, nil
}

// GetFileHashByHash retrieves a file hash by its hash value
func (db *DB) GetFileHashByHash(hash string) (*models.FileHash, error) {
	query := `
		SELECT id, hash, file_size, media_type, created_at
		FROM file_hashes
		WHERE hash = ?
	`

	var fileHash models.FileHash
	err := db.conn.QueryRow(query, hash).Scan(
		&fileHash.ID,
		&fileHash.Hash,
		&fileHash.FileSize,
		&fileHash.MediaType,
		&fileHash.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file hash not found")
		}
		return nil, fmt.Errorf("failed to get file hash: %w", err)
	}

	return &fileHash, nil
}

// GetFileHashByID retrieves a file hash by its ID
func (db *DB) GetFileHashByID(id int64) (*models.FileHash, error) {
	query := `
		SELECT id, hash, file_size, media_type, created_at
		FROM file_hashes
		WHERE id = ?
	`

	var fileHash models.FileHash
	err := db.conn.QueryRow(query, id).Scan(
		&fileHash.ID,
		&fileHash.Hash,
		&fileHash.FileSize,
		&fileHash.MediaType,
		&fileHash.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file hash not found")
		}
		return nil, fmt.Errorf("failed to get file hash: %w", err)
	}

	return &fileHash, nil
}
