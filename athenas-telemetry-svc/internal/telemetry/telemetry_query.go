package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Service.Query is the internal domain logic for querying telemetry data.
//
// It will return an error if the deviceID is empty or if no telemetry data exists for the given
// deviceID.
func (rSvc *Service) Query(deviceID string) ([]models.Telemetry, error) {
	if deviceID == "" {
		return nil, ErrMissingDeviceID
	}

	results, err := rSvc.dataStore.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	return results, nil
}
