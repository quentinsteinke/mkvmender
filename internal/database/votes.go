package database

import (
	"database/sql"
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// CreateOrUpdateVote creates a new vote or updates existing one
func (db *DB) CreateOrUpdateVote(submissionID, userID int64, voteType models.VoteType) error {
	// Check if vote already exists
	existing, err := db.GetVoteBySubmissionAndUser(submissionID, userID)
	if err == nil && existing != nil {
		// Update existing vote
		return db.UpdateVote(existing.ID, voteType)
	}

	// Create new vote
	query := `
		INSERT INTO votes (submission_id, user_id, vote_type)
		VALUES (?, ?, ?)
	`

	_, err = db.conn.Exec(query, submissionID, userID, int(voteType))
	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

// UpdateVote updates an existing vote
func (db *DB) UpdateVote(voteID int64, voteType models.VoteType) error {
	query := `
		UPDATE votes
		SET vote_type = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, int(voteType), voteID)
	if err != nil {
		return fmt.Errorf("failed to update vote: %w", err)
	}

	return nil
}

// GetVoteBySubmissionAndUser retrieves a vote by submission and user
func (db *DB) GetVoteBySubmissionAndUser(submissionID, userID int64) (*models.Vote, error) {
	query := `
		SELECT id, submission_id, user_id, vote_type, created_at, updated_at
		FROM votes
		WHERE submission_id = ? AND user_id = ?
	`

	var vote models.Vote
	var voteTypeInt int

	err := db.conn.QueryRow(query, submissionID, userID).Scan(
		&vote.ID,
		&vote.SubmissionID,
		&vote.UserID,
		&voteTypeInt,
		&vote.CreatedAt,
		&vote.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vote not found")
		}
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	vote.VoteType = models.VoteType(voteTypeInt)

	return &vote, nil
}

// DeleteVote deletes a vote
func (db *DB) DeleteVote(submissionID, userID int64) error {
	query := `
		DELETE FROM votes
		WHERE submission_id = ? AND user_id = ?
	`

	result, err := db.conn.Exec(query, submissionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

// GetVoteCountBySubmission gets vote statistics for a submission
func (db *DB) GetVoteCountBySubmission(submissionID int64) (upvotes, downvotes, score int, err error) {
	query := `
		SELECT
			COUNT(CASE WHEN vote_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN vote_type = -1 THEN 1 END) as downvotes,
			COALESCE(SUM(vote_type), 0) as score
		FROM votes
		WHERE submission_id = ?
	`

	err = db.conn.QueryRow(query, submissionID).Scan(&upvotes, &downvotes, &score)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get vote count: %w", err)
	}

	return upvotes, downvotes, score, nil
}
