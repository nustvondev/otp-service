package main

import (
	"log"

	"github.com/nustvondev/otp-service/config"
	"github.com/nustvondev/otp-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
