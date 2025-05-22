package api

import (
	"net/http"

	"github.com/artnikel/taskscheduler/constants"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/artnikel/taskscheduler/tasks"
	"github.com/labstack/echo/v4"
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

func (h *Handler) CreatePingTask(c echo.Context) error {
	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	taskID := h.Scheduler.AddTask(tasks.MakePingTask(req.Address))
	return c.JSON(http.StatusCreated, map[string]string{"task_id": taskID})
}

func (h *Handler) GetTaskStatus(c echo.Context) error {
	id := c.Param("id")
	task, ok := h.Scheduler.GetTask(id)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	resp := map[string]interface{}{
		"id":     task.ID,
		"status": task.Status,
	}
	if task.Status == constants.StatusDone {
		resp["result"] = task.Result
	}
	if task.Status == constants.StatusFailed {
		resp["error"] = task.Err.Error()
	}

	return c.JSON(http.StatusOK, resp)
}
