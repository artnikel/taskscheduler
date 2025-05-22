package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/artnikel/taskscheduler/api"
	"github.com/artnikel/taskscheduler/config"
	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/artnikel/taskscheduler/tasks"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logging.Error.Fatalf("failed to load config: %v", err)
	}

	err = logging.Init(cfg.Logging.Path)
	if err != nil {
		logging.Error.Fatalf("failed to init logger: %v", err)
	}

	sched := scheduler.NewScheduler(cfg.Scheduler.MaxConcurrentTasks)
	handler := api.NewHandler(sched)

	http.HandleFunc("/tasks/ping", handler.CreatePingTask)
	http.HandleFunc("/tasks/", handler.GetTaskStatus)
	http.HandleFunc("/tasks/stats", handler.GetStats)

	go func() { // worker for server load
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for _, site := range cfg.Worker.PingSites {
				sched.AddTask(tasks.MakePingTask(site))
			}
		}
	}()

	logging.Info.Printf("Server started at :%d\n", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		logging.Error.Fatalf("server error: %v", err)
	}
}
