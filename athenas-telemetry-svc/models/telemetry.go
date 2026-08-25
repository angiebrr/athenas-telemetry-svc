package models

import "errors"

// ================================================================================================

var (
	ErrMissingDeviceID   = errors.New("missing device ID")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
	ErrMissingMetrics    = errors.New("missing metrics")
	ErrMissingMetricName = errors.New("missing metric name")
)

// ------------------------------------------------------------------------------------------------

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

func (rTelemetry *Telemetry) Validate() error {
	if rTelemetry.DeviceID == "" {
		return ErrMissingDeviceID
	}

	if rTelemetry.Timestamp < 0 {
		return ErrInvalidTimestamp
	}

	if len(rTelemetry.Metrics) == 0 {
		return ErrMissingMetrics
	}

	for _, metric := range rTelemetry.Metrics {
		if metric.Name == "" {
			return ErrMissingMetricName
		}
	}

	return nil
}
