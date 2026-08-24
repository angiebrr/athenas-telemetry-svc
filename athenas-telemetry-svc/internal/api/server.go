package api

import "github.com/gin-gonic/gin"

// ================================================================================================

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	router := gin.Default()
	InitHandlers(router)

	return &Server{
		router,
	}
}
