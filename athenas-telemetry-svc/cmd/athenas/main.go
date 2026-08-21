package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================================================================================================

func main() {
	fmt.Println("Telemetry service is running...")

	// TODO: move this code somewhere else
	// TODO: save the telemetry data, define the struct, etc. (this is just a dummy)
	router := gin.Default()
	router.POST("/v1/telemetry", func(ctx *gin.Context) {
		ctx.Status(http.StatusAccepted)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
