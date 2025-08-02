package routes

import (
	"net/http"
	"strings"

	"lwnra-devo-api/handlers"
)

// Router handles HTTP routing for the API
type Router struct {
	devotionalHandler *handlers.DevotionalHandler
	systemHandler     *handlers.SystemHandler
}

// NewRouter creates a new router with handlers
func NewRouter(devotionalHandler *handlers.DevotionalHandler, systemHandler *handlers.SystemHandler) *Router {
	return &Router{
		devotionalHandler: devotionalHandler,
		systemHandler:     systemHandler,
	}
}

// ServeHTTP implements the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route requests based on path and method
	path := strings.TrimSuffix(r.URL.Path, "/")
	
	switch {
	case path == "/api/devotionals" && r.Method == http.MethodGet:
		router.devotionalHandler.GetDevotionals(w, r)
	case strings.HasPrefix(path, "/api/devotionals/") && r.Method == http.MethodGet:
		router.devotionalHandler.GetDevotionalByDate(w, r)
	case path == "/api/devotionals/sync" && r.Method == http.MethodPost:
		router.devotionalHandler.SyncDevotionals(w, r)
	case path == "/api/devotionals/parse" && r.Method == http.MethodPost:
		router.devotionalHandler.ParseDevotional(w, r)
	case path == "/api/scheduler/status" && r.Method == http.MethodGet:
		router.systemHandler.GetSchedulerStatus(w, r)
	case path == "/health" && r.Method == http.MethodGet:
		router.systemHandler.HealthCheck(w, r)
	case path == "/" || path == "":
		router.apiInfo(w, r)
	default:
		router.notFound(w, r)
	}
}

// API info endpoint
func (router *Router) apiInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	apiInfo := `{
		"name": "LWNRA Devotional API",
		"version": "1.0.0",
		"description": "REST API for managing daily devotionals with automated scheduling",
		"endpoints": {
			"GET /api/devotionals": "Get all devotionals (with optional ?limit=N)",
			"GET /api/devotionals/{date}": "Get devotional by date (YYYY-MM-DD format)",
			"POST /api/devotionals/sync": "Sync devotionals from Facebook",
			"POST /api/devotionals/parse": "Parse devotional text",
			"GET /api/scheduler/status": "Get scheduler status and next run time",
			"GET /health": "Health check"
		},
		"scheduler": {
			"sync_time": "4:45 AM Philippine Time (UTC+8)",
			"backup_sync": "5:15 AM Philippine Time (UTC+8)",
			"timezone": "Asia/Manila"
		}
	}`
	w.Write([]byte(apiInfo))
}

// 404 handler
func (router *Router) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"success":false,"error":"Endpoint not found"}`))
}
