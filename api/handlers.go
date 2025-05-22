package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/artnikel/taskscheduler/tasks"
)

type Handler struct {
	Scheduler *scheduler.Scheduler
}

func NewHandler(s *scheduler.Scheduler) *Handler {
	return &Handler{Scheduler: s}
}

type CreateTaskRequest struct {
	Address string `json:"address"`
}

func (h *Handler) CreatePingTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logging.Error.Println("method not allowed")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		logging.Error.Println("invalid request body:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id := h.Scheduler.AddTask(tasks.MakePingTask(req.Address))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"task_id": id})
}

func (h *Handler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		logging.Error.Println("missing task ID in request")
		http.Error(w, "missing task ID", http.StatusBadRequest)
		return
	}
	task, ok := h.Scheduler.GetTask(id)
	if !ok {
		logging.Error.Println("task not found for ID:", id)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"id":     task.ID,
		"status": task.Status,
	}
	if task.Result != "" {
		logging.Error.Println("Task", id, "task result is nil")
		resp["result"] = task.Result
	}
	if task.Err != nil {
		logging.Error.Println("Task", id, "failed with error:", task.Err)
		resp["error"] = task.Err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.Scheduler.GetStats()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stats)
}
