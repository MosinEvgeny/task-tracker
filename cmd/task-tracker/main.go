package main

import (
	"log"

	"github.com/MosinEvgeny/task-tracker/internal/app"
	"github.com/MosinEvgeny/task-tracker/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
