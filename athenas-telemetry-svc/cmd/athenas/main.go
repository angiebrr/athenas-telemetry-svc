// Package main is the entry point for the telemetry service. It initializes and runs the API server.
package main

import (
	"fmt"
	"log"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
)

// ================================================================================================

func main() {
	// TODO: Use postgres data store at some point
	dataStore := data.NewInMemoryStore()
	telemetrySvc := telemetry.NewService(dataStore)

	// TODO: Add server config for port and other things- listens on 0.0.0.0:8080 by default
	server := api.NewServer(telemetrySvc)

	fmt.Println("Telemetry service is running...")
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to run telemetry service: %v", err)
	}
}
