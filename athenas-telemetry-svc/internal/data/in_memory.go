package data

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

type InMemoryDataStore struct {
}

func NewInMemoryDataStore() TelemetryDataStorer {
	return &InMemoryDataStore{}
}

func (rStore *InMemoryDataStore) Insert(data models.Telemetry) error {
	return nil
}

func (rStore *InMemoryDataStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	return nil, nil
}
