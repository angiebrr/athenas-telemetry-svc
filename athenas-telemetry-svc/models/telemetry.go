package models

// ================================================================================================

// TODO: Move bindings / validations so it's done on the business layer instead of the transport
// layer

// Telemetry is a model representing ingested metrics sent to the telemetry service.
type Telemetry struct {
	DeviceID  string          `json:"device_id" binding:"required"`
	Timestamp int64           `json:"timestamp" binding:"required,min=0"`
	Metrics   []MetricReading `json:"metrics" binding:"required"`
}

// MetricReading is an individual metric sent to the service.
//
// NOTE: It is designed this way so devs can send telemetry without having to update this model
// but still allowing for easy/cheap deserialization.
type MetricReading struct {
	Name  string  `json:"name" binding:"required"`
	Value float64 `json:"value" binding:"required"`
}
