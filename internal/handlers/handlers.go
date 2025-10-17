package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/quentinsteinke/mkvmender/internal/database"
	"github.com/quentinsteinke/mkvmender/internal/models"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	db *database.DB
}

// New creates a new Handler
func New(db *database.DB) *Handler {
	return &Handler{db: db}
}

// RegisterHandler handles user registration
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" {
		respondError(w, http.StatusBadRequest, "username is required")
		return
	}

	// Create user
	user, err := h.db.CreateUser(req.Username)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// LookupHandler handles file hash lookup
func (h *Handler) LookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get hash from query parameter
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash parameter is required")
		return
	}

	// Get file hash info
	fileHash, err := h.db.GetFileHashByHash(hash)
	if err != nil {
		respondJSON(w, http.StatusOK, models.HashLookupResponse{
			Hash:        hash,
			Submissions: []models.SubmissionWithVotes{},
		})
		return
	}

	// Get submissions for this hash
	submissions, err := h.db.GetSubmissionsByHash(hash)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get submissions")
		return
	}

	// Get metadata for each submission
	for i := range submissions {
		meta, err := h.db.GetMetadataBySubmissionID(submissions[i].ID)
		if err == nil && meta != nil {
			submissions[i].Metadata = meta
		}
	}

	response := models.HashLookupResponse{
		Hash:        fileHash.Hash,
		FileSize:    fileHash.FileSize,
		MediaType:   fileHash.MediaType,
		Submissions: submissions,
	}

	respondJSON(w, http.StatusOK, response)
}

// UploadHandler handles naming submission upload
func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get authenticated user
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req models.UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if req.Hash == "" || req.Filename == "" {
		respondError(w, http.StatusBadRequest, "hash and filename are required")
		return
	}

	if req.MediaType != models.MediaTypeMovie && req.MediaType != models.MediaTypeTV {
		respondError(w, http.StatusBadRequest, "invalid media type")
		return
	}

	// Create or get file hash
	fileHash, err := h.db.CreateFileHash(req.Hash, req.FileSize, req.MediaType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create file hash")
		return
	}

	// Create submission
	submission, err := h.db.CreateSubmission(fileHash.ID, user.ID, req.Filename)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create submission")
		return
	}

	// Create metadata if provided
	if req.Metadata != nil {
		req.Metadata.SubmissionID = submission.ID
		if err := h.db.CreateMetadata(req.Metadata); err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	respondJSON(w, http.StatusCreated, submission)
}

// VoteHandler handles voting on submissions
func (h *Handler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get authenticated user
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate vote type
	if req.VoteType != models.VoteUp && req.VoteType != models.VoteDown {
		respondError(w, http.StatusBadRequest, "invalid vote type")
		return
	}

	// Create or update vote
	if err := h.db.CreateOrUpdateVote(req.SubmissionID, user.ID, req.VoteType); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create vote")
		return
	}

	// Get updated submission with vote counts
	submission, err := h.db.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get submission")
		return
	}

	// Return updated vote counts
	response := map[string]interface{}{
		"success":    true,
		"message":    "vote recorded",
		"upvotes":    submission.Upvotes,
		"downvotes":  submission.Downvotes,
		"vote_score": submission.VoteScore,
	}

	respondJSON(w, http.StatusOK, response)
}

// DeleteVoteHandler handles removing a vote
func (h *Handler) DeleteVoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get authenticated user
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Get submission ID from query parameter
	submissionIDStr := r.URL.Query().Get("submission_id")
	if submissionIDStr == "" {
		respondError(w, http.StatusBadRequest, "submission_id parameter is required")
		return
	}

	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid submission_id")
		return
	}

	// Delete vote
	if err := h.db.DeleteVote(submissionID, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete vote")
		return
	}

	respondSuccess(w, "vote removed")
}

// SearchHandler handles searching by title
func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get search query from parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	// Get sort parameter (default: relevance)
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "relevance"
	}

	var sortBy database.SortBy
	switch sortParam {
	case "relevance":
		sortBy = database.SortByRelevance
	case "votes":
		sortBy = database.SortByVotes
	case "date":
		sortBy = database.SortByDate
	case "title":
		sortBy = database.SortByTitle
	default:
		sortBy = database.SortByRelevance
	}

	// Get fuzzy parameter (default: true)
	useFuzzy := r.URL.Query().Get("fuzzy") != "false"

	// Search database
	dbResults, err := h.db.SearchByTitle(query, sortBy, useFuzzy)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "search failed")
		return
	}

	// Convert database results to API results
	var results []models.SearchResult
	for _, dbResult := range dbResults {
		result := models.SearchResult{
			Title:       dbResult.Title,
			Year:        dbResult.Year,
			MediaType:   dbResult.MediaType,
			Season:      dbResult.Season,
			Episode:     dbResult.Episode,
			Hash:        dbResult.Hash,
			FileSize:    dbResult.FileSize,
			Submissions: dbResult.Submissions,
		}
		results = append(results, result)
	}

	response := models.SearchResponse{
		Query:   query,
		Results: results,
	}

	respondJSON(w, http.StatusOK, response)
}

// HealthHandler handles health check
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
