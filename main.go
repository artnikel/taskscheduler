package main

import (
	"log"

	"github.com/artnikel/taskscheduler/api"
	"github.com/artnikel/taskscheduler/internal/logging"
	"github.com/artnikel/taskscheduler/scheduler"
	"github.com/labstack/echo/v4"
)

func main() {
	err := logging.Init("logs")
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	logging.Info.Println("Starting application")
	sched := scheduler.NewScheduler(3)
	e := echo.New()

	handler := api.NewHandler(sched)

	e.POST("/tasks/ping", handler.CreatePingTask)
	e.GET("/tasks/:id", handler.GetTaskStatus)

	e.Logger.Fatal(e.Start(":8080"))
}
