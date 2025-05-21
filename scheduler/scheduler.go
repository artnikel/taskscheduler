package scheduler

import (
	"sync"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusRunning TaskStatus = "running"
	StatusDone    TaskStatus = "done"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID     string
	Status TaskStatus
	Result string
	Err    error
}

type TaskFunc func() (string, error)

type Scheduler struct {
	maxConcurrent int
	tasks         map[string]*Task
	taskLock      sync.RWMutex
	sem           chan struct{}
}

func NewScheduler(maxConcurrent int) *Scheduler {
	return &Scheduler{
		maxConcurrent: maxConcurrent,
		tasks:         make(map[string]*Task),
		sem:           make(chan struct{}, maxConcurrent),
	}
}

func (s *Scheduler) AddTask(fn TaskFunc) string {
	taskID := uuid.NewString()
	task := &Task{
		ID:     taskID,
		Status: StatusPending,
	}

	s.taskLock.Lock()
	s.tasks[taskID] = task
	s.taskLock.Unlock()

	go s.runTask(taskID, fn)
	return taskID
}

func (s *Scheduler) runTask(taskID string, fn TaskFunc) {
	s.sem <- struct{}{} 
	defer func() { <-s.sem }()

	s.taskLock.Lock()
	task, exists := s.tasks[taskID]
	if !exists {
		s.taskLock.Unlock()
		return
	}
	task.Status = StatusRunning
	s.taskLock.Unlock()

	result, err := fn()

	s.taskLock.Lock()
	defer s.taskLock.Unlock()

	if err != nil {
		task.Status = StatusFailed
		task.Err = err
		return
	}
	task.Status = StatusDone
	task.Result = result
}

func (s *Scheduler) GetTask(id string) (*Task, bool) {
	s.taskLock.RLock()
	defer s.taskLock.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}
