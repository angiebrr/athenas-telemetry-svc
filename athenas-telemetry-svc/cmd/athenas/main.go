package main

import (
	"fmt"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
)

// ================================================================================================

func main() {
	server := api.NewServer()

	fmt.Println("Telemetry service is running...")
	server.Run() // TODO: Add config for port and other things- listens on 0.0.0.0:8080 by default
}
