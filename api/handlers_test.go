package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/artnikel/taskscheduler/constants"
	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
)

func NewLoggerForTest() *logging.Logger {
	return &logging.Logger{
		Info:  log.New(io.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(io.Discard, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
func TestCreatePingTask_Valid(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	body := []byte(`{"address": "example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	h.CreatePingTask(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var data map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&data)
	if data["task_id"] == "" {
		t.Fatal("task_id not returned")
	}
}

func TestCreatePingTask_Invalid(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte(`{}`)))
	w := httptest.NewRecorder()
	h.CreatePingTask(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetTaskStatus_Valid(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	id := s.AddTask(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "pong", nil
	})
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+id, http.NoBody)
	w := httptest.NewRecorder()
	h.GetTaskStatus(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var data map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&data)
	if data["status"] != string(constants.StatusDone) {
		t.Errorf("expected status 'done', got %v", data["status"])
	}
}

func TestGetTaskStatus_NotFound(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	req := httptest.NewRequest(http.MethodGet, "/tasks/nonexistent", http.NoBody)
	w := httptest.NewRecorder()
	h.GetTaskStatus(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetStats(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	_ = s.AddTask(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "result", nil
	})
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/stats", http.NoBody)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCreateStatusTask_Valid(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	body := []byte(`{"url": "http://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/status-tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateStatusTask(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", resp.StatusCode)
	}
	var data map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatal("failed to decode response:", err)
	}
	if data["task_id"] == "" {
		t.Fatal("task_id not returned")
	}
}

func TestCreateStatusTask_InvalidMethod(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	req := httptest.NewRequest(http.MethodGet, "/status-tasks", http.NoBody)
	w := httptest.NewRecorder()

	h.CreateStatusTask(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405 Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestCreateStatusTask_InvalidBody(t *testing.T) {
	s := scheduler.NewScheduler(1)
	logger := NewLoggerForTest()
	h := NewHandler(s, logger)

	body := []byte(`{"url": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/status-tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateStatusTask(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}
