// Package main is the entry point for the telemetry service. It initializes and runs the API server.
package main

import (
	"fmt"
	"log"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
	"github.com/angiebrr/athenas-telemetry-svc/internal/config"
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
)

// ================================================================================================

func main() {
	// Get env vars / .env files and throw them in a Config struct
	conf, err := config.InitEnv()
	if err != nil {
		log.Fatalf("failed to init env: %v", err)
	}

	// Set up the backend service that manages the telemetry
	// TODO: Use postgres data store at some point
	dataStore := data.NewInMemoryStore()
	telemetrySvc := telemetry.NewService(dataStore)

	// Set up the HTTP server that serves the telemetry from the backend service
	server := api.NewServer(telemetrySvc)

	fmt.Println("Telemetry service is running...")
	err = server.Run(fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatalf("failed to run telemetry service: %v", err)
	}
}
