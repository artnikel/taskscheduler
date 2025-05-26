package tasks

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeGetStatusTask_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	task := MakeGetStatusTask(server.URL)
	result, err := task()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestMakeGetStatusTask_FailureStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	task := MakeGetStatusTask(server.URL)
	result, err := task()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
}

func TestMakeGetStatusTask_ConnectionError(t *testing.T) {
	task := MakeGetStatusTask("http://invalid.localhost")

	result, err := task()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
}
