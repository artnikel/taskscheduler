package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/artnikel/taskscheduler/constants"
	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
)

func init() {
	_ = logging.Init("testlogs") 
}

func TestMain(m *testing.M) {
	_ = logging.Init("testlogs")
	os.Exit(m.Run())
}

func TestCreatePingTask_Valid(t *testing.T) {
	s := scheduler.NewScheduler(1)
	h := NewHandler(s)

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
	h := NewHandler(s)

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
	h := NewHandler(s)

	id := s.AddTask(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "pong", nil
	})
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+id, nil)
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
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/tasks/nonexistent", nil)
	w := httptest.NewRecorder()
	h.GetTaskStatus(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetStats(t *testing.T) {
	s := scheduler.NewScheduler(1)
	h := NewHandler(s)

	_ = s.AddTask(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "result", nil
	})
	time.Sleep(50 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
