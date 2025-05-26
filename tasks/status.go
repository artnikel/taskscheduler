package tasks

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/artnikel/taskscheduler/constants"
)

// MakeGetStatusTask returns a task that sends an HTTP GET request to the given URL.
func MakeGetStatusTask(url string) func() (string, error) {
	return func() (string, error) {
		client := &http.Client{
			Timeout: constants.TaskTimeout,
		}
		start := time.Now()
		resp, err := client.Get(url)
		elapsed := time.Since(start)

		if err != nil {
			return "", fmt.Errorf("http get %s failed: %w", url, err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Println("error closing response body:", err)
			}
		}()

		if resp.StatusCode >= http.StatusBadRequest {
			return "", fmt.Errorf("http get %s returned error status: %d", url, resp.StatusCode)
		}

		return fmt.Sprintf("http get %s success, status: %d, time: %v", url, resp.StatusCode, elapsed), nil
	}
}
