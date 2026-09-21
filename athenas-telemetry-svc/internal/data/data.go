// Package data defines the interface and errors for the data layer of the telemetry service.
package data

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

// NotFoundError is returned when a query for telemetry data does not return any results.
type NotFoundError struct {
	DeviceID string
}

// NotFoundError.Error returns a string representation of the error, including the deviceID
// for which no telemetry data was found.
func (rErr NotFoundError) Error() string {
	return "data not found for device: " + rErr.DeviceID
}

// ------------------------------------------------------------------------------------------------

// TelemetryDataStorer is an interface that defines the methods for storing and retrieving telemetry data.
type TelemetryDataStorer interface {
	// Insert adds new telemetry data to the store.
	Insert(data ...models.Telemetry) error
	// GetByDeviceID retrieves telemetry data for a given deviceID from the store.
	// If no telemetry data exists for the given deviceID, it will return a NotFoundError.
	GetByDeviceID(deviceID string) ([]models.Telemetry, error)
}
