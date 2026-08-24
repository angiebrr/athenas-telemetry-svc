package api

import "github.com/gin-gonic/gin"

// ================================================================================================

// Server encapsulates the main gin.Engine instance and any other configuration needed to host and
// server the telemetry service
type Server struct {
	*gin.Engine
}

// NewServer creates a Server instance by setting up the gin.Engine instance and initializing its
// handlers
func NewServer() *Server {
	router := gin.Default()
	InitHandlers(router)

	return &Server{
		router,
	}
}
