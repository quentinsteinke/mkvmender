package database

import (
	"database/sql"
	"fmt"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// AdminListSubmissions retrieves a paginated list of submissions with filters
func (db *DB) AdminListSubmissions(page, limit int, userID *int64, sortBy string) ([]models.AdminSubmissionListItem, int, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if userID != nil {
		whereClause += " AND ns.user_id = ?"
		args = append(args, *userID)
	}

	// Build ORDER BY clause
	orderBy := "ORDER BY ns.created_at DESC"
	switch sortBy {
	case "votes":
		orderBy = "ORDER BY vote_score DESC, ns.created_at DESC"
	case "title":
		orderBy = "ORDER BY nm.title ASC, ns.created_at DESC"
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM naming_submissions ns
		%s
	`, whereClause)

	var total int
	err := db.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Get submissions with pagination
	query := fmt.Sprintf(`
		SELECT
			ns.id,
			ns.hash_id,
			ns.user_id,
			ns.filename,
			ns.created_at,
			fh.hash,
			fh.file_size,
			fh.media_type,
			u.username,
			u.role,
			COALESCE(SUM(v.vote_type), 0) as vote_score,
			COUNT(CASE WHEN v.vote_type = 1 THEN 1 END) as upvotes,
			COUNT(CASE WHEN v.vote_type = -1 THEN 1 END) as downvotes
		FROM naming_submissions ns
		JOIN file_hashes fh ON ns.hash_id = fh.id
		JOIN users u ON ns.user_id = u.id
		LEFT JOIN votes v ON ns.id = v.submission_id
		LEFT JOIN naming_metadata nm ON ns.id = nm.submission_id
		%s
		GROUP BY ns.id, ns.hash_id, ns.user_id, ns.filename, ns.created_at,
		         fh.hash, fh.file_size, fh.media_type, u.username, u.role
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, limit, offset)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []models.AdminSubmissionListItem
	for rows.Next() {
		var item models.AdminSubmissionListItem
		err := rows.Scan(
			&item.ID,
			&item.HashID,
			&item.UserID,
			&item.Filename,
			&item.CreatedAt,
			&item.Hash,
			&item.FileSize,
			&item.MediaType,
			&item.Username,
			&item.UserRole,
			&item.VoteScore,
			&item.Upvotes,
			&item.Downvotes,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return submissions, total, nil
}

// AdminListUsers retrieves a paginated list of users with filters
func (db *DB) AdminListUsers(page, limit int, role, status string) ([]models.AdminUserListItem, int, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if role != "" {
		whereClause += " AND u.role = ?"
		args = append(args, role)
	}

	if status == "active" {
		whereClause += " AND u.is_active = 1"
	} else if status == "suspended" {
		whereClause += " AND u.is_active = 0"
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM users u
		%s
	`, whereClause)

	var total int
	err := db.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with submission counts
	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.username,
			u.role,
			u.is_active,
			COUNT(ns.id) as submission_count,
			u.created_at
		FROM users u
		LEFT JOIN naming_submissions ns ON u.id = ns.user_id
		%s
		GROUP BY u.id, u.username, u.role, u.is_active, u.created_at
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.AdminUserListItem
	for rows.Next() {
		var user models.AdminUserListItem
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.IsActive,
			&user.SubmissionCount,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return users, total, nil
}

// DeleteSubmission deletes a submission and all associated data
func (db *DB) DeleteSubmission(submissionID int64) error {
	// Foreign keys with CASCADE will handle deletion of votes and metadata
	query := `DELETE FROM naming_submissions WHERE id = ?`

	result, err := db.conn.Exec(query, submissionID)
	if err != nil {
		return fmt.Errorf("failed to delete submission: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("submission not found")
	}

	return nil
}

// LogModerationAction logs an admin action
func (db *DB) LogModerationAction(adminID int64, actionType, targetType string, targetID int64, reason string) error {
	query := `
		INSERT INTO moderation_actions (admin_id, action_type, target_type, target_id, reason)
		VALUES (?, ?, ?, ?, ?)
	`

	var reasonPtr *string
	if reason != "" {
		reasonPtr = &reason
	}

	_, err := db.conn.Exec(query, adminID, actionType, targetType, targetID, reasonPtr)
	if err != nil {
		return fmt.Errorf("failed to log moderation action: %w", err)
	}

	return nil
}

// GetAdminStats retrieves system statistics for the admin dashboard
func (db *DB) GetAdminStats() (*models.AdminStats, error) {
	stats := &models.AdminStats{}

	// Get total users
	err := db.conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get active users
	err = db.conn.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = 1").Scan(&stats.ActiveUsers)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Get total submissions
	err = db.conn.QueryRow("SELECT COUNT(*) FROM naming_submissions").Scan(&stats.TotalSubmissions)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count submissions: %w", err)
	}

	// Get total votes
	err = db.conn.QueryRow("SELECT COUNT(*) FROM votes").Scan(&stats.TotalVotes)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count votes: %w", err)
	}

	// Pending actions is always 0 for now (no approval workflow yet)
	stats.PendingActions = 0

	return stats, nil
}
