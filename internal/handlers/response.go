package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// respondSuccess sends a JSON success response
func respondSuccess(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: message,
	})
}
