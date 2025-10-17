package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/quentinsteinke/mkvmender/internal/models"
)

// Client represents an API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a new API client
func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("%s: %s", errResp.Error, errResp.Message)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Register registers a new user
func (c *Client) Register(username string) (*models.User, error) {
	req := map[string]string{"username": username}
	var user models.User
	if err := c.doRequest("POST", "/api/register", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Lookup looks up naming submissions by hash
func (c *Client) Lookup(hash string) (*models.HashLookupResponse, error) {
	path := fmt.Sprintf("/api/lookup?hash=%s", hash)
	var response models.HashLookupResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Upload uploads a new naming submission
func (c *Client) Upload(req *models.UploadRequest) (*models.NamingSubmission, error) {
	var submission models.NamingSubmission
	if err := c.doRequest("POST", "/api/upload", req, &submission); err != nil {
		return nil, err
	}
	return &submission, nil
}

// Vote votes on a submission
func (c *Client) Vote(submissionID int64, voteType models.VoteType) error {
	req := models.VoteRequest{
		SubmissionID: submissionID,
		VoteType:     voteType,
	}
	return c.doRequest("POST", "/api/vote", req, nil)
}

// DeleteVote removes a vote
func (c *Client) DeleteVote(submissionID int64) error {
	path := fmt.Sprintf("/api/vote/delete?submission_id=%d", submissionID)
	return c.doRequest("DELETE", path, nil, nil)
}

// Search searches for naming submissions by title
func (c *Client) Search(query, sortBy string, useFuzzy bool) (*models.SearchResponse, error) {
	params := url.Values{}
	params.Add("q", query)
	if sortBy != "" {
		params.Add("sort", sortBy)
	}
	if !useFuzzy {
		params.Add("fuzzy", "false")
	}

	path := fmt.Sprintf("/api/search?%s", params.Encode())
	var response models.SearchResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Health checks the API health
func (c *Client) Health() error {
	var result map[string]string
	return c.doRequest("GET", "/api/health", nil, &result)
}
