// Package main is an entry point to application
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/artnikel/taskscheduler/api"
	"github.com/artnikel/taskscheduler/config"
	"github.com/artnikel/taskscheduler/constants"
	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/artnikel/taskscheduler/tasks"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := logging.NewLogger(cfg.Logging.Path)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sched := scheduler.NewScheduler(cfg.Scheduler.MaxConcurrentTasks)
	handler := api.NewHandler(sched, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/ping", handler.CreatePingTask)
	mux.HandleFunc("/tasks/", handler.GetTaskStatus)
	mux.HandleFunc("/tasks/stats", handler.GetStats)

	go func() { // worker for server load
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for _, site := range cfg.Worker.PingSites {
				sched.AddTask(tasks.MakePingTask(site))
			}
		}
	}()

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  constants.ServerTimeout,
		WriteTimeout: constants.ServerTimeout,
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), constants.ServerTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error.Fatalf("http server shutdown error %v", err)
		}
		close(stopped)
	}()

	logger.Info.Printf("starting HTTP server on :%d\n", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error.Fatalf("http server not listening: %v", err)
	}

	<-stopped
}
