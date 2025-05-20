package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusQueued  TaskStatus = "queued"
	StatusRunning TaskStatus = "running"
	StatusDone    TaskStatus = "done"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID       string
	Status   TaskStatus
	Result   string
	Err      error
	doneChan chan struct{}
	work     func() (string, error)
}

type Scheduler struct {
	tasks       map[string]*Task
	queue       []string
	maxParallel int
	semaphore   chan struct{}
	mu          sync.RWMutex
}

var errorLogger *log.Logger

func init() {
	f, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open error log file: %v", err)
	}
	errorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewScheduler(maxParallel int) *Scheduler {
	return &Scheduler{
		tasks:       make(map[string]*Task),
		queue:       make([]string, 0),
		maxParallel: maxParallel,
		semaphore:   make(chan struct{}, maxParallel),
	}
}

func MakePingTask(address string) func() (string, error) {
	return func() (string, error) {
		timeout := 2 * time.Second
		start := time.Now()

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, "80"), timeout)
		elapsed := time.Since(start)

		if err != nil {
			return "", fmt.Errorf("ping %s failed: %v", address, err)
		}

		conn.Close()
		return fmt.Sprintf("ping %s success, time: %v", address, elapsed), nil
	}
}

func (s *Scheduler) AddTask(work func() (string, error)) string {
	id := uuid.New().String()
	task := &Task{
		ID:       id,
		Status:   StatusQueued,
		work:     work,
		doneChan: make(chan struct{}),
	}

	s.mu.Lock()
	s.tasks[id] = task
	s.queue = append(s.queue, id)
	s.mu.Unlock()

	go s.runTask(id)
	return id
}

func (s *Scheduler) runTask(taskID string) {
	fmt.Printf("Task %s is running...\n", taskID)

	s.semaphore <- struct{}{}

	s.mu.Lock()
	task, exists := s.tasks[taskID]
	if !exists {
		s.mu.Unlock()
		<-s.semaphore
		return
	}
	task.Status = StatusRunning
	s.mu.Unlock()

	result, err := task.work()

	s.mu.Lock()
	if err != nil {
		task.Status = StatusFailed
		task.Err = err
		errorLogger.Printf("Task %s failed: %v", taskID, err)
	} else {
		task.Status = StatusDone
		task.Result = result
	}
	close(task.doneChan)
	s.mu.Unlock()

	<-s.semaphore
	fmt.Printf("Task %s is finished with status: %s\n", taskID, task.Status)

}

func (s *Scheduler) GetStatus(taskID string) (TaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[taskID]
	if !exists {
		return "", errors.New("task not found")
	}
	return task.Status, nil
}

func (s *Scheduler) GetResult(taskID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[taskID]
	if !exists {
		return "", errors.New("task not found")
	}
	if task.Status != StatusDone {
		return "", errors.New("task not completed yet")
	}
	return task.Result, nil
}

func main() {
    scheduler := NewScheduler(3)

    targets := []string{
        "google.com",
        "yandex.com",
        "youtube.com",
    }

    for _, host := range targets {
        id := scheduler.AddTask(MakePingTask(host))
        fmt.Println("Ping task added with ID:", id)
    }

    time.Sleep(5 * time.Second)
}

