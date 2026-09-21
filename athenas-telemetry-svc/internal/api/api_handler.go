package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Path constants for the telemetry service API endpoints.
const (
	// TelemetryPath represents the telemetry ingestion resource HTTP path.
	// This is the path that clients will POST telemetry data to.
	TelemetryPath = "/v1/telemetry"

	// QueryTelemetryBasePath represents the base path for querying telemetry data, excluding the
	// device ID path parameter.
	QueryTelemetryBasePath = TelemetryPath

	// DeviceIDPathName is the name of the path parameter used to specify the device ID when querying
	// telemetry data.
	DeviceIDPathName = "device_id"

	// QueryTelemetryPath represents the full path for querying telemetry data by device ID.
	QueryTelemetryPath = QueryTelemetryBasePath + "/:" + DeviceIDPathName
)

// ------------------------------------------------------------------------------------------------

// TODO: Add the ingester/querier as a dependency to make the dependency explicit by passing it in
// to the handler funcs

// InitHandlers registers HTTP endpoint handlers for the telemetry service.
func InitHandlers(router *gin.Engine, dataStore data.Storer) {
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
func HandleIngestTelemetry(ctx *gin.Context, dataStore data.Storer) {
	// deserialize request into a telemetry model
	var newData models.Telemetry
	if err := ctx.ShouldBindJSON(&newData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := telemetry.Ingest(newData, dataStore); err != nil {
		var statusCode int
		errMsg := err.Error()
		if _, isErr := errors.AsType[telemetry.ValidateTelemetryError](err); isErr {
			// if the error is a validation error, return a 400 Bad Request
			statusCode = http.StatusBadRequest
		} else {
			// otherwise, just return a 500 Internal Server Error and make the error vague
			statusCode = http.StatusInternalServerError
			errMsg = "internal server error"
			slog.Warn("failed to ingest telemetry", "error", err)
		}
		ctx.JSON(statusCode, gin.H{"error": errMsg})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// HandleQueryTelemetry handles requests to query telemetry data for a specific device ID.
func HandleQueryTelemetry(ctx *gin.Context, dataStore data.Storer) {
	deviceID := ctx.Param(DeviceIDPathName)

	queriedData, err := telemetry.Query(deviceID, dataStore)
	if err != nil {
		var statusCode int
		errMsg := err.Error()
		if _, isErr := errors.AsType[telemetry.ValidateTelemetryError](err); isErr {
			// if the error is a validation error, return a 400 Bad Request
			statusCode = http.StatusBadRequest
		} else if _, isErr := errors.AsType[data.NotFoundError](err); isErr {
			// if we couldn't find the data, return a 404 Not Found
			statusCode = http.StatusNotFound
		} else {
			// otherwise, just return a 500 Internal Server Error and make the error vague
			statusCode = http.StatusInternalServerError
			errMsg = "internal server error"
			slog.Warn("failed to query telemetry", "error", err)
		}
		ctx.JSON(statusCode, gin.H{"error": errMsg})
		return
	}

	ctx.JSON(http.StatusOK, queriedData)
}
