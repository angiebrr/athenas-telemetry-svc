package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

func main() {
	fmt.Println("Telemetry service is running...")

	// TODO: move this code somewhere else
	// TODO: save the telemetry data, define the struct, etc. (this is just a dummy)
	router := gin.Default()
	router.POST("/v1/telemetry", func(ctx *gin.Context) {
		var data specifications.Telemetry

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: do something with data

		ctx.Status(http.StatusAccepted)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
