// Package models provides the data models used in the application
package models

import "github.com/artnikel/taskscheduler/constants"

// Task entity
type Task struct {
	ID     string
	Status constants.TaskStatus
	Result string
	Err    error
}
