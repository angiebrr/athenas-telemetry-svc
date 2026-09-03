package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// ValidateTelemetryError is a custom error that is used to tell the difference between a validation
// error and other internal system errors, which is useful for things like HTTP response codes.
type ValidateTelemetryError struct {
	err string
}

// ValidateTelemetryError.Error returns the validation message this error was built with.
func (rErr ValidateTelemetryError) Error() string {
	return rErr.err
}

var (
	// ErrMissingDeviceID is returned when the telemetry data is missing a device ID
	ErrMissingDeviceID = ValidateTelemetryError{err: "missing device ID"}
	// ErrInvalidTimestamp is returned when the telemetry data has an invalid timestamp (negative)
	ErrInvalidTimestamp = ValidateTelemetryError{err: "invalid timestamp"}
	// ErrMissingMetrics is returned when the telemetry data is missing metrics.
	ErrMissingMetrics = ValidateTelemetryError{err: "missing metrics"}
	// ErrMissingMetricName is returned when the telemetry data has a metric with a missing name
	ErrMissingMetricName = ValidateTelemetryError{err: "missing metric name"}
)

// ------------------------------------------------------------------------------------------------

// Validate verifies that the telemetry data is valid and returns an error if it is not.
//
// This is intentionally NOT on the model (which is public facing) to make sure our logic is properly
// contained and our acceptance tests aren't able to bypass the validation logic by calling the
// model directly.
func Validate(data models.Telemetry) error {
	if data.DeviceID == "" {
		return ErrMissingDeviceID
	}

	if data.Timestamp < 0 {
		return ErrInvalidTimestamp
	}

	if len(data.Metrics) == 0 {
		return ErrMissingMetrics
	}

	for _, metric := range data.Metrics {
		if metric.Name == "" {
			return ErrMissingMetricName
		}
	}

	return nil
}
