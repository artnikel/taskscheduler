package tasks

import (
	"fmt"
	"net"
	"time"
)

func MakePingTask(address string) func() (string, error) {
	return func() (string, error) {
		timeout := 2 * time.Second
		start := time.Now()
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, "80"), timeout)
		elapsed := time.Since(start)
		if err != nil {
			return "", fmt.Errorf("ping %s failed: %w", address, err)
		}
		_ = conn.Close()
		return fmt.Sprintf("ping %s success, time: %v", address, elapsed), nil
	}
}
