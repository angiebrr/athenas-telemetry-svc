package data

import (
	"slices"
	"sync"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// InMemoryStore is a simple in-memory implementation of the data.Storer interface. It is intended
// for use in testing and development scenarios where a persistent data store is not required.
//
// Although it is not for production use, it is thread-safe.
type InMemoryStore struct {
	// storedData is a map of deviceID to telemetry data for that device.
	storedData map[string][]models.Telemetry

	// mu is a mutex to ensure thread-safe access to the storedData map.
	mu sync.Mutex
}

// NewInMemoryStore creates a new instance of InMemoryStore, ready for data insertions
// and queries.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		storedData: make(map[string][]models.Telemetry),
	}
}

// Insert adds new telemetry data to the in-memory store.
//
// It is thread-safe.
func (rStore *InMemoryStore) Insert(newTelemetry ...models.Telemetry) error {
	rStore.mu.Lock()
	defer rStore.mu.Unlock()

	for _, currData := range newTelemetry {
		rStore.storedData[currData.DeviceID] = append(rStore.storedData[currData.DeviceID], currData)
	}

	return nil
}

// GetByDeviceID retrieves telemetry data for a given deviceID from the in-memory store.
//
// If no telemetry data exists for the given deviceID, it will return a NotFoundError.
//
// The returned results are a copy so that the caller can safely read them without worrying
// about concurrent modifications to the underlying data.
func (rStore *InMemoryStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	rStore.mu.Lock()
	defer rStore.mu.Unlock()

	// retrieve telemetry from device
	results, exists := rStore.storedData[deviceID]
	if !exists {
		return nil, NotFoundError{DeviceID: deviceID}
	}

	// copy over results to make sure the caller will have thread-safe reads
	//
	// NOTE: can't just use slices.Clone(results) because we need a deepy-copy
	clonedResults := make([]models.Telemetry, len(results))
	for idx, currData := range results {
		clonedResults[idx] = currData
		clonedResults[idx].Metrics = slices.Clone(currData.Metrics)
	}

	return clonedResults, nil
}
