// Package tasks provides utilities for creating executable tasks
package tasks

import (
	"fmt"
	"net"
	"time"

	"github.com/artnikel/taskscheduler/constants"
)

// MakePingTask returns a task function that pings the given address over TCP
func MakePingTask(address string) func() (string, error) {
	return func() (string, error) {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, "80"), constants.TaskTimeout)
		elapsed := time.Since(start)
		if err != nil {
			return "", fmt.Errorf("ping %s failed: %w", address, err)
		}
		_ = conn.Close()
		return fmt.Sprintf("ping %s success, time: %v", address, elapsed), nil
	}
}
