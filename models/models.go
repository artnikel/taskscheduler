package models

import "github.com/artnikel/taskscheduler/constants"

type Task struct {
	ID     string
	Status constants.TaskStatus
	Result string
	Err    error
}