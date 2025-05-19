package main

import (
    "fmt"
    "sync"
    "time"
    "errors"
    "github.com/google/uuid"
)

type TaskStatus string

const (
    StatusQueued   TaskStatus = "queued"
    StatusRunning  TaskStatus = "running"
    StatusDone     TaskStatus = "done"
    StatusFailed   TaskStatus = "failed"
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

func NewScheduler(maxParallel int) *Scheduler {
    return &Scheduler{
        tasks:       make(map[string]*Task),
        queue:       make([]string, 0),
        maxParallel: maxParallel,
        semaphore:   make(chan struct{}, maxParallel),
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
    scheduler := NewScheduler(2)

    for i := 1; i <= 5; i++ {
        id := scheduler.AddTask(func() (string, error) {
            time.Sleep(2 * time.Second)
            return fmt.Sprintf("Task completed at %v", time.Now()), nil
        })
        fmt.Println("Task added with ID:", id)
    }

    time.Sleep(10 * time.Second) 
}
