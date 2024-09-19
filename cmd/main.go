package main

import (
	"log"

	"github.com/FlyKarlik/auth-service/internal/app"
	"github.com/FlyKarlik/auth-service/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config initialization error: %s", err)
	}

	app := app.New(cfg)
	if app.Run(); err != nil {
		log.Fatalf("Application run error: %s", err)
	}
}
