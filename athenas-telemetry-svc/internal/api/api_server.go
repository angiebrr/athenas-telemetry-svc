// Package api provides the API server implementation for the telemetry service, including the
// gin.Engine setup and handler initialization.
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
)

// ================================================================================================

// Server encapsulates the main gin.Engine instance and any other configuration needed to host and
// serve the telemetry service
type Server struct {
	*gin.Engine
	telemetrySvc *telemetry.Service // The internal domain logic service for managing telemetry
}

// NewServer creates a Server instance by setting up the gin.Engine instance and initializing its
// handlers
func NewServer(telemetrySvc *telemetry.Service) *Server {
	router := gin.Default()
	InitHandlers(router, telemetrySvc)

	return &Server{
		router,
		telemetrySvc,
	}
}
