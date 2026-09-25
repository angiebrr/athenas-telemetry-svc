// Package shared contains shared utilities and helpers for the telemetry service, including test utilities.
package shared

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// ValidTelemetry generates a valid telemetry data object for testing purposes.
//
// It creates a random device ID as a UUID, a current unix timestamp, and a single metric reading
// with a random temperature value between 0 and 120.
func ValidTelemetry() models.Telemetry {
	const maxTemp = 120.0
	randTemp := rand.Float64() * maxTemp

	data := models.Telemetry{
		DeviceID:  uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: randTemp},
		},
	}

	return data
}

// ValidTelemetryN generates a slice of n valid telemetry data objects for testing purposes.
//
// See ValidTelemetry for details on how each telemetry object is generated.
func ValidTelemetryN(n int) []models.Telemetry {
	telemetryList := make([]models.Telemetry, n)
	for range n {
		telemetryList = append(telemetryList, ValidTelemetry())
	}
	return telemetryList
}
