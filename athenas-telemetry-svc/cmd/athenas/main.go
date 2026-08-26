// Package main is the entry point for the telemetry service. It initializes and runs the API server.
package main

import (
	"fmt"
	"log"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
)

// ================================================================================================

func main() {
	// TODO: Add server config for port and other things- listens on 0.0.0.0:8080 by default
	server := api.NewServer()

	fmt.Println("Telemetry service is running...")
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to run telemetry service: %v", err)
	}
}
