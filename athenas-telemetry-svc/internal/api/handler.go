package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

const (
	deviceIDName = "device_id"
)

const (
	// TelemetryPath represents the telemetry ingestion resource HTTP path.
	// This is the path that clients will POST telemetry data to.
	TelemetryPath = "/v1/telemetry"

	// QueryTelemetryPath represents the telemetry query resource HTTP path.
	// This is the path that clients will GET telemetry data from.
	QueryTelemetryPath = TelemetryPath + "/:" + deviceIDName
)

// ------------------------------------------------------------------------------------------------

// TODO: Add the ingester/querier as a dependency to make the dependency explicit by passing it in
// to the handler funcs

// InitHandlers registers HTTP endpoint handlers for the telemetry service.
func InitHandlers(router *gin.Engine, dataStore data.TelemetryDataStorer) {
	router.POST(TelemetryPath, func(ctx *gin.Context) {
		HandleIngestTelemetry(ctx, dataStore)
	})
	router.GET(QueryTelemetryPath, func(ctx *gin.Context) {
		HandleQueryTelemetry(ctx, dataStore)
	})
}

// ------------------------------------------------------------------------------------------------

// HandleIngestTelemetry binds and validates telemetry sent to the ingest endpoint, answering
// 202 Accepted once the domain has taken it.
//
// NOTE: The payload is not persisted or dispatched anywhere yet -- see the TODO in ingest.Ingest.
func HandleIngestTelemetry(ctx *gin.Context, dataStore data.TelemetryDataStorer) {
	// deserialize request into a telemetry model
	var newData models.Telemetry
	if err := ctx.ShouldBindJSON(&newData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := telemetry.Ingest(newData, dataStore); err != nil {
		// if the error is a validation error, return a 400 Bad Request
		// otherwise, return a 500 Internal Server Error
		statusCode := http.StatusInternalServerError
		if _, ok := errors.AsType[telemetry.ValidateTelemetryError](err); ok {
			statusCode = http.StatusBadRequest
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusAccepted)
}

func HandleQueryTelemetry(ctx *gin.Context, dataStore data.TelemetryDataStorer) {
	deviceID := ctx.Param(deviceIDName)

	queriedData, err := telemetry.Query(deviceID, dataStore) // TODO: Account for device ID validation
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, queriedData)
}
