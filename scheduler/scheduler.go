// Package scheduler manages task execution with concurrency control
package scheduler

import (
	"sync"

	"github.com/artnikel/taskscheduler/constants"
	"github.com/artnikel/taskscheduler/models"
	"github.com/google/uuid"
)

// TaskFunc defines the function signature for a scheduled task
type TaskFunc func() (string, error)

// Scheduler handles task management and concurrent execution
type Scheduler struct {
	maxConcurrent int
	tasks         map[string]*models.Task
	taskLock      sync.RWMutex
	sem           chan struct{}
}

// NewScheduler creates a new Scheduler with the given concurrency limit
func NewScheduler(maxConcurrent int) *Scheduler {
	return &Scheduler{
		maxConcurrent: maxConcurrent,
		tasks:         make(map[string]*models.Task),
		sem:           make(chan struct{}, maxConcurrent),
	}
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
	task.Status = constants.StatusRunning
	s.taskLock.Unlock()

	result, err := fn()

	s.taskLock.Lock()
	defer s.taskLock.Unlock()

	if err != nil {
		task.Status = constants.StatusFailed
		task.Err = err
		return
	}
	task.Status = constants.StatusDone
	task.Result = result
}

// AddTask adds a new task to the scheduler and runs it asynchronously
func (s *Scheduler) AddTask(fn TaskFunc) string {
	taskID := uuid.NewString()
	task := &models.Task{
		ID:     taskID,
		Status: constants.StatusPending,
	}

	s.taskLock.Lock()
	s.tasks[taskID] = task
	s.taskLock.Unlock()

	go s.runTask(taskID, fn)
	return taskID
}

// GetTask returns the task with the given ID, if it exists
func (s *Scheduler) GetTask(id string) (*models.Task, bool) {
	s.taskLock.RLock()
	defer s.taskLock.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

// GetStats returns the count of tasks by their status
func (s *Scheduler) GetStats() map[constants.TaskStatus]int {
	stats := map[constants.TaskStatus]int{
		constants.StatusPending: 0,
		constants.StatusRunning: 0,
		constants.StatusDone:    0,
		constants.StatusFailed:  0,
	}

	s.taskLock.RLock()
	defer s.taskLock.RUnlock()

	for _, task := range s.tasks {
		stats[task.Status]++
	}

	return stats
}
