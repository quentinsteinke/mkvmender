package models

import "time"

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeTV    MediaType = "tv"
)

// VoteType represents upvote or downvote
type VoteType int

const (
	VoteDown VoteType = -1
	VoteUp   VoteType = 1
)

// UserRole represents user permission levels
type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	APIKey    string    `json:"api_key,omitempty"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileHash represents a hashed media file
type FileHash struct {
	ID        int64     `json:"id"`
	Hash      string    `json:"hash"`
	FileSize  int64     `json:"file_size"`
	MediaType MediaType `json:"media_type"`
	CreatedAt time.Time `json:"created_at"`
}

// NamingSubmission represents a user's submission for a file name
type NamingSubmission struct {
	ID        int64     `json:"id"`
	HashID    int64     `json:"hash_id"`
	UserID    int64     `json:"user_id"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vote represents a user's vote on a naming submission
type Vote struct {
	ID           int64     `json:"id"`
	SubmissionID int64     `json:"submission_id"`
	UserID       int64     `json:"user_id"`
	VoteType     VoteType  `json:"vote_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NamingMetadata represents additional metadata for a naming submission
type NamingMetadata struct {
	ID           int64     `json:"id"`
	SubmissionID int64     `json:"submission_id"`
	Title        *string   `json:"title,omitempty"`
	Year         *int      `json:"year,omitempty"`
	Season       *int      `json:"season,omitempty"`
	Episode      *int      `json:"episode,omitempty"`
	Quality      *string   `json:"quality,omitempty"`
	Source       *string   `json:"source,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// SubmissionWithVotes represents a naming submission with vote counts
type SubmissionWithVotes struct {
	ID         int64     `json:"id"`
	HashID     int64     `json:"hash_id"`
	UserID     int64     `json:"user_id"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	Hash       string    `json:"hash"`
	FileSize   int64     `json:"file_size"`
	MediaType  MediaType `json:"media_type"`
	Username   string    `json:"username"`
	VoteScore  int       `json:"vote_score"`
	Upvotes    int       `json:"upvotes"`
	Downvotes  int       `json:"downvotes"`
	Metadata   *NamingMetadata `json:"metadata,omitempty"`
}

// HashLookupRequest represents a request to lookup naming submissions by hash
type HashLookupRequest struct {
	Hash     string `json:"hash"`
	FileSize int64  `json:"file_size"`
}

// HashLookupResponse represents the response with available naming options
type HashLookupResponse struct {
	Hash        string                 `json:"hash"`
	FileSize    int64                  `json:"file_size"`
	MediaType   MediaType              `json:"media_type"`
	Submissions []SubmissionWithVotes  `json:"submissions"`
}

// UploadRequest represents a request to upload a new naming submission
type UploadRequest struct {
	Hash      string             `json:"hash"`
	FileSize  int64              `json:"file_size"`
	MediaType MediaType          `json:"media_type"`
	Filename  string             `json:"filename"`
	Metadata  *NamingMetadata    `json:"metadata,omitempty"`
}

// VoteRequest represents a request to vote on a submission
type VoteRequest struct {
	SubmissionID int64    `json:"submission_id"`
	VoteType     VoteType `json:"vote_type"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// SearchRequest represents a search request by title
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Title       string                `json:"title"`
	Year        *int                  `json:"year,omitempty"`
	MediaType   MediaType             `json:"media_type"`
	Season      *int                  `json:"season,omitempty"`
	Episode     *int                  `json:"episode,omitempty"`
	Hash        string                `json:"hash"`
	FileSize    int64                 `json:"file_size"`
	Submissions []SubmissionWithVotes `json:"submissions"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

// ModerationAction represents an admin action
type ModerationAction struct {
	ID         int64     `json:"id"`
	AdminID    int64     `json:"admin_id"`
	ActionType string    `json:"action_type"`
	TargetType string    `json:"target_type"`
	TargetID   int64     `json:"target_id"`
	Reason     *string   `json:"reason,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// AdminStats represents system statistics for admin dashboard
type AdminStats struct {
	TotalUsers       int `json:"total_users"`
	ActiveUsers      int `json:"active_users"`
	TotalSubmissions int `json:"total_submissions"`
	TotalVotes       int `json:"total_votes"`
	PendingActions   int `json:"pending_actions"`
}

// AdminSubmissionListItem represents a submission in the admin list
type AdminSubmissionListItem struct {
	SubmissionWithVotes
	UserRole UserRole `json:"user_role"`
}

// AdminUserListItem represents a user in the admin list
type AdminUserListItem struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Role             UserRole  `json:"role"`
	IsActive         bool      `json:"is_active"`
	SubmissionCount  int       `json:"submission_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// ChangeRoleRequest represents a request to change user role
type ChangeRoleRequest struct {
	Role UserRole `json:"role"`
}

// ChangeStatusRequest represents a request to change user active status
type ChangeStatusRequest struct {
	IsActive bool    `json:"is_active"`
	Reason   *string `json:"reason,omitempty"`
}

// DeleteSubmissionRequest represents a request to delete a submission
type DeleteSubmissionRequest struct {
	Reason *string `json:"reason,omitempty"`
}
