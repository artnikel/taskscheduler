package tasks

import (
	"strings"
	"testing"
)

func TestMakePingTask_Success(t *testing.T) {
	task := MakePingTask("google.com")
	result, err := task()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(result, "ping google.com success") {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMakePingTask_Failure(t *testing.T) {
	task := MakePingTask("nonexistent.domain.local")
	_, err := task()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "ping nonexistent.domain.local failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}
