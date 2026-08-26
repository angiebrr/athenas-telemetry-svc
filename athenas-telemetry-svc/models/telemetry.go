// Package models contains the data models used by the telemetry service, and is exported publicly
// mostly so it can be used by the acceptance tests (but has the benefit of being public for other
// clients as well)
package models

// ================================================================================================

// Telemetry is a model representing ingested metrics sent to the telemetry service.
type Telemetry struct {
	DeviceID  string          `json:"device_id"`
	Timestamp int64           `json:"timestamp"`
	Metrics   []MetricReading `json:"metrics"`
}

// MetricReading is an individual metric sent to the service.
//
// NOTE: It is designed this way so devs can send telemetry without having to update this model
// but still allowing for easy/cheap deserialization.
type MetricReading struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
