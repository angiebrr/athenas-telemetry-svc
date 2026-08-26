// Package main is the entry point for the telemetry service. It initializes and runs the API server.
package main

import (
	"fmt"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
)

// ================================================================================================

func main() {
	server := api.NewServer()

	fmt.Println("Telemetry service is running...")
	err := server.Run() // TODO: Add config for port and other things- listens on 0.0.0.0:8080 by default
	if err != nil {
		panic(fmt.Sprintf("Error running server: %v\n", err))
	}
}
