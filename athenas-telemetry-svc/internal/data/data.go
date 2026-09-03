package data

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

type DataNotFoundError struct {
	DeviceID string
}

func (rErr DataNotFoundError) Error() string {
	return "data not found for device: " + rErr.DeviceID
}

// ------------------------------------------------------------------------------------------------

type TelemetryDataStorer interface {
	Insert(data ...models.Telemetry) error
	GetByDeviceID(deviceID string) ([]models.Telemetry, error)
}
