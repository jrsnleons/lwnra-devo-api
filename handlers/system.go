package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"lwnra-devo-api/scheduler"
)

// SystemHandler handles system-related endpoints
type SystemHandler struct {
	scheduler *scheduler.Scheduler
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(sched *scheduler.Scheduler) *SystemHandler {
	return &SystemHandler{
		scheduler: sched,
	}
}

// SchedulerStatusResponse represents the scheduler status response
type SchedulerStatusResponse struct {
	IsRunning bool      `json:"is_running"`
	NextRun   time.Time `json:"next_run,omitempty"`
	Timezone  string    `json:"timezone"`
}

// GetSchedulerStatus handles GET /api/scheduler/status
func (h *SystemHandler) GetSchedulerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nextRun time.Time
	isRunning := h.scheduler.IsRunning()
	
	if isRunning {
		nextRun = h.scheduler.GetNextRun()
	}

	response := APIResponse{
		Success: true,
		Data: SchedulerStatusResponse{
			IsRunning: isRunning,
			NextRun:   nextRun,
			Timezone:  "Asia/Manila",
		},
	}

	json.NewEncoder(w).Encode(response)
}

// HealthCheck handles GET /health
func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: map[string]interface{}{
			"timestamp": time.Now(),
			"version":   "1.0.0",
		},
	}

	json.NewEncoder(w).Encode(response)
}
