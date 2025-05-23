package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/artnikel/taskscheduler/constants"
)

func TestAddTask_Success(t *testing.T) {
	s := NewScheduler(2)

	id := s.AddTask(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "ok", nil
	})

	time.Sleep(200 * time.Millisecond)
	task, ok := s.GetTask(id)
	if !ok {
		t.Fatal("task should exist")
	}
	if task.Status != constants.StatusDone {
		t.Errorf("expected status %s, got %s", constants.StatusDone, task.Status)
	}
	if task.Result != "ok" {
		t.Errorf("expected result 'ok', got %s", task.Result)
	}
}

func TestAddTask_Failure(t *testing.T) {
	s := NewScheduler(1)

	id := s.AddTask(func() (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "", fmt.Errorf("failed")
	})

	time.Sleep(100 * time.Millisecond)
	task, ok := s.GetTask(id)
	if !ok {
		t.Fatal("task should exist")
	}
	if task.Status != constants.StatusFailed {
		t.Errorf("expected status %s, got %s", constants.StatusFailed, task.Status)
	}
	if task.Err == nil || task.Err.Error() != "failed" {
		t.Errorf("expected error 'failed', got %v", task.Err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	s := NewScheduler(1)
	_, ok := s.GetTask("nonexistent")
	if ok {
		t.Error("expected false, got true")
	}
}
