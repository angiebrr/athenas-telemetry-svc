package data

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

type DataNotFoundError struct {
	DeviceID string
}

func (rErr *DataNotFoundError) Error() string {
	return "data not found for device: " + rErr.DeviceID
}

// ------------------------------------------------------------------------------------------------

type InMemoryDataStore struct {
	data map[string][]models.Telemetry
}

func NewInMemoryDataStore() TelemetryDataStorer {
	return &InMemoryDataStore{
		data: make(map[string][]models.Telemetry),
	}
}

func (rStore *InMemoryDataStore) Insert(data models.Telemetry) error {
	rStore.data[data.DeviceID] = append(rStore.data[data.DeviceID], data)
	return nil
}

func (rStore *InMemoryDataStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	results, exists := rStore.data[deviceID]
	if !exists {
		return nil, &DataNotFoundError{DeviceID: deviceID}
	}
	return results, nil
}
