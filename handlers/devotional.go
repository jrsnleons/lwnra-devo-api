package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
	"lwnra-devo-api/parser"
)

// DevotionalHandler handles all devotional-related API endpoints
type DevotionalHandler struct {
	db       *database.DB
	fbClient *facebook.Client
}

// NewDevotionalHandler creates a new devotional handler
func NewDevotionalHandler(db *database.DB, fbClient *facebook.Client) *DevotionalHandler {
	return &DevotionalHandler{
		db:       db,
		fbClient: fbClient,
	}
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GetDevotionals handles GET /api/devotionals
func (h *DevotionalHandler) GetDevotionals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	devotionals, err := h.db.GetDevotionals(limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch devotionals", err)
		return
	}

	respondWithSuccess(w, "Devotionals retrieved successfully", devotionals)
}

// GetDevotionalByDate handles GET /api/devotionals/{date}
func (h *DevotionalHandler) GetDevotionalByDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract date from URL path
	date := extractDateFromPath(r.URL.Path)
	if date == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", nil)
		return
	}

	devotional, err := h.db.GetDevotionalByDate(date)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Devotional not found for the specified date", err)
		return
	}

	respondWithSuccess(w, "Devotional retrieved successfully", devotional)
}

// SyncDevotionals handles POST /api/devotionals/sync
func (h *DevotionalHandler) SyncDevotionals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch posts from Facebook
	posts, err := h.fbClient.GetRecentPosts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch Facebook posts", err)
		return
	}

	// Filter for devotional posts
	devotionalPosts := facebook.FilterDevotionalPosts(posts)

	// Process and save devotionals
	count := 0
	var errors []string

	for _, post := range devotionalPosts {
		devo := parser.ParseDevotional(post.Message)

		// Use post date if devotional date is empty
		if devo.Date == "" {
			postDate := extractAndFormatPostDate(post.CreatedTime)
			if postDate != "" {
				devo.Date = postDate
			} else {
				devo.Date = time.Now().Format("January 2, 2006")
			}
		}

		// Save to database
		if err := h.db.SaveDevotional(devo); err != nil {
			errors = append(errors, "Failed to save devotional '"+devo.Title+"': "+err.Error())
		} else {
			count++
		}
	}

	response := map[string]interface{}{
		"synced_count": count,
		"total_posts":  len(devotionalPosts),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	respondWithSuccess(w, "Sync completed", response)
}

// ParseDevotional handles POST /api/devotionals/parse
func (h *DevotionalHandler) ParseDevotional(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON request body", err)
		return
	}

	if request.Message == "" {
		respondWithError(w, http.StatusBadRequest, "Message field is required", nil)
		return
	}

	// Parse the devotional
	devo := parser.ParseDevotional(request.Message)

	respondWithSuccess(w, "Devotional parsed successfully", devo)
}

// Helper functions
func respondWithSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	response := APIResponse{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func extractDateFromPath(path string) string {
	// Extract date from path like "/api/devotionals/2025-08-02"
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3] // date part
	}
	return ""
}

func extractAndFormatPostDate(createdTime string) string {
	t, err := time.Parse("2006-01-02T15:04:05-0700", createdTime)
	if err != nil {
		return ""
	}
	return t.Format("January 2, 2006")
}
