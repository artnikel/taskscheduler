// Package api contains HTTP handlers for task management
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/artnikel/taskscheduler/tasks"
)

// Handler provides HTTP endpoints backed by a Scheduler
type Handler struct {
	Scheduler *scheduler.Scheduler
	Logger    *logging.Logger
}

// NewHandler creates a new Handler with the given Scheduler
func NewHandler(s *scheduler.Scheduler, logger *logging.Logger) *Handler {
	return &Handler{Scheduler: s, Logger: logger}
}

// CreateTaskRequest represents a request to create a ping task
type CreateTaskRequest struct {
	Address string `json:"address"`
}

// CreatePingTask handles POST requests to add a new ping task
func (h *Handler) CreatePingTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.Error.Println("method not allowed")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		h.Logger.Error.Println("invalid request body:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id := h.Scheduler.AddTask(tasks.MakePingTask(req.Address))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"task_id": id})
}

// GetTaskStatus handles GET requests to retrieve task status by ID
func (h *Handler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		h.Logger.Error.Println("missing task ID in request")
		http.Error(w, "missing task ID", http.StatusBadRequest)
		return
	}
	task, ok := h.Scheduler.GetTask(id)
	if !ok {
		h.Logger.Error.Println("task not found for ID:", id)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"id":     task.ID,
		"status": task.Status,
	}
	if task.Result != "" {
		h.Logger.Error.Println("Task", id, "task result is nil")
		resp["result"] = task.Result
	}
	if task.Err != nil {
		h.Logger.Error.Println("Task", id, "failed with error:", task.Err)
		resp["error"] = task.Err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetStats handles GET requests to retrieve aggregated task statistics
func (h *Handler) GetStats(w http.ResponseWriter, _ *http.Request) {
	stats := h.Scheduler.GetStats()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}
