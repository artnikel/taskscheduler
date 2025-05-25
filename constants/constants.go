// Package constants defines shared constants used across the application
package constants

import "time"

// TaskStatus represents the state of a task
type TaskStatus string

const (
	// StatusPending - Task is queued for execution
	StatusPending TaskStatus = "pending"
	// StatusRunning - Task is currently executing
	StatusRunning TaskStatus = "running"
	// StatusDone - Task completed successfully
	StatusDone TaskStatus = "done"
	// StatusFailed - Task execution failed
	StatusFailed TaskStatus = "failed"
	// TaskTimeout - Maximum allowed time for task execution
	TaskTimeout = 2 * time.Second
	// ServerTimeout is read and write timeout of server config
	ServerTimeout = 10 * time.Second
	// DirPerm - Directory permission
	DirPerm = 0o750
	// FilePerm - File permission
	FilePerm = 0o600
)
