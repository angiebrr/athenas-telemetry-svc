package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

func Query(deviceID string, dataStore data.TelemetryDataStorer) ([]models.Telemetry, error) {
	if deviceID == "" {
		return nil, ErrQueryMissingDeviceID
	}

	data, err := dataStore.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
