package data

import (
	"slices"
	"sync"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

type InMemoryDataStore struct {
	data map[string][]models.Telemetry
	mu   sync.Mutex
}

func NewInMemoryDataStore() TelemetryDataStorer {
	return &InMemoryDataStore{
		data: make(map[string][]models.Telemetry),
	}
}

func (rStore *InMemoryDataStore) Insert(newTelemetry ...models.Telemetry) error {
	rStore.mu.Lock()
	defer rStore.mu.Unlock()

	for _, currData := range newTelemetry {
		rStore.data[currData.DeviceID] = append(rStore.data[currData.DeviceID], currData)
	}

	return nil
}

func (rStore *InMemoryDataStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	rStore.mu.Lock()
	defer rStore.mu.Unlock()

	// retrieve telemetry from device
	results, exists := rStore.data[deviceID]
	if !exists {
		return nil, DataNotFoundError{DeviceID: deviceID}
	}

	// copy over results to make sure the caller will have thread-safe reads
	newResults := slices.Clone(results)

	return newResults, nil
}
