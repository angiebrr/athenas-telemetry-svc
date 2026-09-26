package telemetry_test

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

// telemetryWithDeviceID copies the given data but with the given device ID.
func telemetryWithDeviceID(data models.Telemetry, deviceID string) models.Telemetry {
	data.DeviceID = deviceID
	return data
}

// telemetryWithTimestamp copies the given data but with the given timestamp.
func telemetryWithTimestamp(data models.Telemetry, timestamp int64) models.Telemetry {
	data.Timestamp = timestamp
	return data
}

// telemetryWithMetrics copies the given data but with the given metrics.
func telemetryWithMetrics(data models.Telemetry, metrics []models.MetricReading) models.Telemetry {
	data.Metrics = metrics
	return data
}
