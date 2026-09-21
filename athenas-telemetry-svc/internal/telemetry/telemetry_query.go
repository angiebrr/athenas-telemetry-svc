package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Query is the internal domain logic for querying telemetry data.
//
// It will return an error if the deviceID is empty or if no telemetry data exists for the given
// deviceID.
func Query(deviceID string, dataStore data.Storer) ([]models.Telemetry, error) {
	if deviceID == "" {
		return nil, ErrMissingDeviceID
	}

	results, err := dataStore.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	return results, nil
}
