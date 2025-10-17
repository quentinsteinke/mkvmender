package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/quentinsteinke/mkvmender/internal/database"
	"github.com/quentinsteinke/mkvmender/internal/models"
)

// AdminHandler holds dependencies for admin HTTP handlers
type AdminHandler struct {
	db *database.DB
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// ListSubmissionsHandler handles listing all submissions with pagination and filters
// GET /api/admin/submissions?page=1&limit=50&sort=date&user_id=123
func (h *AdminHandler) ListSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Parse filters
	userIDStr := r.URL.Query().Get("user_id")
	var userID *int64
	if userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			userID = &id
		}
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "date"
	}

	// Fetch submissions from database
	submissions, total, err := h.db.AdminListSubmissions(page, limit, userID, sortBy)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch submissions")
		return
	}

	response := map[string]interface{}{
		"submissions": submissions,
		"total":       total,
		"page":        page,
		"limit":       limit,
	}

	respondJSON(w, http.StatusOK, response)
}

// GetSubmissionHandler handles getting a single submission with full details
// GET /api/admin/submissions/:id
func (h *AdminHandler) GetSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract submission ID from URL path
	submissionIDStr := r.URL.Query().Get("id")
	if submissionIDStr == "" {
		respondError(w, http.StatusBadRequest, "submission ID is required")
		return
	}

	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid submission ID")
		return
	}

	submission, err := h.db.GetSubmissionByID(submissionID)
	if err != nil {
		respondError(w, http.StatusNotFound, "submission not found")
		return
	}

	// Get metadata
	metadata, _ := h.db.GetMetadataBySubmissionID(submissionID)
	if metadata != nil {
		submission.Metadata = metadata
	}

	respondJSON(w, http.StatusOK, submission)
}

// DeleteSubmissionHandler handles deleting a submission
// DELETE /api/admin/submissions/:id
func (h *AdminHandler) DeleteSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get admin user from context
	admin, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Extract submission ID
	submissionIDStr := r.URL.Query().Get("id")
	if submissionIDStr == "" {
		respondError(w, http.StatusBadRequest, "submission ID is required")
		return
	}

	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid submission ID")
		return
	}

	// Parse reason from request body
	var req models.DeleteSubmissionRequest
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&req)
	}

	// Delete submission
	if err := h.db.DeleteSubmission(submissionID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete submission")
		return
	}

	// Log moderation action
	reasonStr := ""
	if req.Reason != nil {
		reasonStr = *req.Reason
	}
	h.db.LogModerationAction(admin.ID, "delete_submission", "submission", submissionID, reasonStr)

	respondSuccess(w, "submission deleted successfully")
}

// ListUsersHandler handles listing all users with pagination and filters
// GET /api/admin/users?page=1&limit=50&role=admin&status=active
func (h *AdminHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Parse filters
	role := r.URL.Query().Get("role")
	status := r.URL.Query().Get("status")

	// Fetch users from database
	users, total, err := h.db.AdminListUsers(page, limit, role, status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	respondJSON(w, http.StatusOK, response)
}

// GetUserHandler handles getting a single user with full details
// GET /api/admin/users/:id
func (h *AdminHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract user ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ChangeUserRoleHandler handles changing a user's role
// PUT /api/admin/users/:id/role
func (h *AdminHandler) ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get admin user from context
	admin, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Extract user ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Parse request body
	var req models.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate role
	if req.Role != models.RoleUser && req.Role != models.RoleModerator && req.Role != models.RoleAdmin {
		respondError(w, http.StatusBadRequest, "invalid role")
		return
	}

	// Prevent admins from changing their own role (security measure)
	if userID == admin.ID {
		respondError(w, http.StatusForbidden, "cannot change your own role")
		return
	}

	// Update user role
	if err := h.db.UpdateUserRole(userID, req.Role); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update user role")
		return
	}

	// Log moderation action
	h.db.LogModerationAction(admin.ID, "change_role", "user", userID, string(req.Role))

	respondSuccess(w, "user role updated successfully")
}

// ChangeUserStatusHandler handles activating or suspending a user
// PUT /api/admin/users/:id/status
func (h *AdminHandler) ChangeUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get admin user from context
	admin, ok := GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Extract user ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Parse request body
	var req models.ChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Prevent admins from suspending themselves
	if userID == admin.ID {
		respondError(w, http.StatusForbidden, "cannot change your own status")
		return
	}

	// Update user status
	if err := h.db.UpdateUserStatus(userID, req.IsActive); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update user status")
		return
	}

	// Log moderation action
	actionType := "suspend_user"
	if req.IsActive {
		actionType = "activate_user"
	}
	reasonStr := ""
	if req.Reason != nil {
		reasonStr = *req.Reason
	}
	h.db.LogModerationAction(admin.ID, actionType, "user", userID, reasonStr)

	message := "user suspended successfully"
	if req.IsActive {
		message = "user activated successfully"
	}
	respondSuccess(w, message)
}

// GetStatsHandler handles fetching admin dashboard statistics
// GET /api/admin/stats
func (h *AdminHandler) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats, err := h.db.GetAdminStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch statistics")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
