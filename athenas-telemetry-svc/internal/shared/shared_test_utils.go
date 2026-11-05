// Package shared contains shared utilities and helpers for the telemetry service, like test and string utils.
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
// It will call ValidTelemetryWithDeviceID with a random device ID as a UUID.
func ValidTelemetry() models.Telemetry {
	deviceID := uuid.New().String()
	return ValidTelemetryWithDeviceID(deviceID)
}

// ValidTelemetryWithDeviceID generates a valid telemetry data object for testing purposes.
//
// It uses the given device ID, a current unix timestamp, and a single metric reading with a
// random temperature value between 0 and 120.
func ValidTelemetryWithDeviceID(deviceID string) models.Telemetry {
	const maxTemp = 120.0
	randTemp := rand.Float64() * maxTemp

	data := models.Telemetry{
		DeviceID:  deviceID,
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: randTemp},
		},
	}

	return data
}

// ValidTelemetryN generates a slice of n valid telemetry data objects for testing purposes.
// NOTE: It will make telemetry data points with the same device ID
//
// See ValidTelemetry and ValidTelemetryWithDeviceID for details on how each telemetry object is
// generated.
func ValidTelemetryN(n int) ([]models.Telemetry, string) {
	if n <= 0 {
		return nil, ""
	}
	telemetryList := make([]models.Telemetry, n)

	// Set up first datapoint so we can extract the device ID
	firstData := ValidTelemetry()
	deviceID := firstData.DeviceID
	telemetryList[0] = firstData

	for idx := 1; idx < n; idx++ {
		telemetryList[idx] = ValidTelemetryWithDeviceID(deviceID)
	}

	return telemetryList, deviceID
}
