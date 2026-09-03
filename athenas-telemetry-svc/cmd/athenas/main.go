// Package main is the entry point for the telemetry service. It initializes and runs the API server.
package main

import (
	"fmt"
	"log"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
)

// ================================================================================================

func main() {
	// TODO: Use postgres data store at some point
	dataStore := data.NewInMemoryDataStore()

	// TODO: Add server config for port and other things- listens on 0.0.0.0:8080 by default
	server := api.NewServer(dataStore)

	fmt.Println("Telemetry service is running...")
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to run telemetry service: %v", err)
	}
}
