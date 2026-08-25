package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/internal/ingest"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

const (
	TelemetryPath = "/v1/telemetry"
)

// ------------------------------------------------------------------------------------------------

// InitHandlers registers HTTP endpoint handlers for the telemetry service.
func InitHandlers(router *gin.Engine) {
	router.POST(TelemetryPath, HandleIngestTelemetry)
}

// ------------------------------------------------------------------------------------------------

// HandleIngestTelemetry handles the endpoint that validates and stores the telemetry that's sent
// to the telemetry service.
func HandleIngestTelemetry(ctx *gin.Context) {
	// deserialize request into a telemetry model
	var data models.Telemetry
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ingest.Ingest(data); err != nil {
		// if the error is a validation error, return a 400 Bad Request
		// otherwise, return a 500 Internal Server Error
		statusCode := http.StatusInternalServerError
		if _, ok := errors.AsType[ingest.ValidateTelemetryError](err); ok {
			statusCode = http.StatusBadRequest
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusAccepted)
}
