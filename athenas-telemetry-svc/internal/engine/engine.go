// Package engine handles managing the ingestion of telemetry data at scale
package engine

import (
	"errors"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Queue represents a data structure that allows callers to insert pending telemetry asynchronously
type Queue struct {
	store data.Storer
}

// NewQueue creates a new Queue instance with the given data.Storer as the backing data store
// that the queue workers will use to manage telemetry.
func NewQueue(store data.Storer) (*Queue, error) {
	if store == nil {
		return nil, errors.New("data store is required for engine.Queue")
	}

	return &Queue{
		store,
	}, nil
}

// ------------------------------------------------------------------------------------------------

// Queue.PushTelemetry
// TODO: Make docstring, it's not worth doing until it's fully implemented
func (rQueue *Queue) PushTelemetry(newTelemetry ...models.Telemetry) error {
	// FIXME: Make queue
	go func() {
		_ = rQueue.store.Insert(newTelemetry...)
	}()
	return nil
}
