package database

import (
	"database/sql"
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// CreateSubmission creates a new naming submission
func (db *DB) CreateSubmission(hashID, userID int64, filename string) (*models.NamingSubmission, error) {
	query := `
		INSERT INTO naming_submissions (hash_id, user_id, filename)
		VALUES (?, ?, ?)
		RETURNING id, hash_id, user_id, filename, created_at, updated_at
	`

	var submission models.NamingSubmission
	err := db.conn.QueryRow(query, hashID, userID, filename).Scan(
		&submission.ID,
		&submission.HashID,
		&submission.UserID,
		&submission.Filename,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	return &submission, nil
}

// GetSubmissionsByHash retrieves all naming submissions for a given hash with vote counts
func (db *DB) GetSubmissionsByHash(hash string) ([]models.SubmissionWithVotes, error) {
	query := `
		SELECT
			id, hash_id, user_id, filename, created_at,
			hash, file_size, media_type, username,
			vote_score, upvotes, downvotes
		FROM submissions_with_votes
		WHERE hash = ?
		ORDER BY vote_score DESC, created_at DESC
	`

	rows, err := db.conn.Query(query, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []models.SubmissionWithVotes
	for rows.Next() {
		var s models.SubmissionWithVotes
		var mediaTypeStr string

		err := rows.Scan(
			&s.ID,
			&s.HashID,
			&s.UserID,
			&s.Filename,
			&s.CreatedAt,
			&s.Hash,
			&s.FileSize,
			&mediaTypeStr,
			&s.Username,
			&s.VoteScore,
			&s.Upvotes,
			&s.Downvotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}

		s.MediaType = models.MediaType(mediaTypeStr)
		submissions = append(submissions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return submissions, nil
}

// GetSubmissionByID retrieves a submission by its ID
func (db *DB) GetSubmissionByID(id int64) (*models.SubmissionWithVotes, error) {
	query := `
		SELECT
			id, hash_id, user_id, filename, created_at,
			hash, file_size, media_type, username,
			vote_score, upvotes, downvotes
		FROM submissions_with_votes
		WHERE id = ?
	`

	var s models.SubmissionWithVotes
	var mediaTypeStr string

	err := db.conn.QueryRow(query, id).Scan(
		&s.ID,
		&s.HashID,
		&s.UserID,
		&s.Filename,
		&s.CreatedAt,
		&s.Hash,
		&s.FileSize,
		&mediaTypeStr,
		&s.Username,
		&s.VoteScore,
		&s.Upvotes,
		&s.Downvotes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	s.MediaType = models.MediaType(mediaTypeStr)

	return &s, nil
}

// CreateMetadata creates naming metadata for a submission
func (db *DB) CreateMetadata(meta *models.NamingMetadata) error {
	query := `
		INSERT INTO naming_metadata (submission_id, title, year, season, episode, quality, source)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		meta.SubmissionID,
		meta.Title,
		meta.Year,
		meta.Season,
		meta.Episode,
		meta.Quality,
		meta.Source,
	)
	if err != nil {
		return fmt.Errorf("failed to create metadata: %w", err)
	}

	return nil
}

// GetMetadataBySubmissionID retrieves metadata for a submission
func (db *DB) GetMetadataBySubmissionID(submissionID int64) (*models.NamingMetadata, error) {
	query := `
		SELECT id, submission_id, title, year, season, episode, quality, source, created_at
		FROM naming_metadata
		WHERE submission_id = ?
	`

	var meta models.NamingMetadata
	err := db.conn.QueryRow(query, submissionID).Scan(
		&meta.ID,
		&meta.SubmissionID,
		&meta.Title,
		&meta.Year,
		&meta.Season,
		&meta.Episode,
		&meta.Quality,
		&meta.Source,
		&meta.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No metadata is OK
		}
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	return &meta, nil
}
