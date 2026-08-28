package data

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

type TelemetryDataStorer interface {
	Insert(data ...models.Telemetry) error
	GetByDeviceID(deviceID string) ([]models.Telemetry, error)
}
