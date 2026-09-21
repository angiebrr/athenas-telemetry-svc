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
	randTemp := rand.Float64() * 120.0 // rand temp from 0. to 120.

	data := models.Telemetry{
		DeviceID:  uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: randTemp},
		},
	}

	return data
}
