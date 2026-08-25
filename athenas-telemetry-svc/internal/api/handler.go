package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// InitHandlers registers HTTP endpoint handlers for the telemetry service.
func InitHandlers(router *gin.Engine) {
	router.POST("/v1/telemetry", HandleIngestTelemetry)
}

// ------------------------------------------------------------------------------------------------

// HandleIngestTelemetry handles the endpoint that validates and stores the telemetry that's sent
// to the telemetry service.
func HandleIngestTelemetry(ctx *gin.Context) {
	var data models.Telemetry
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := data.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: do something with data

	ctx.Status(http.StatusAccepted)
}
