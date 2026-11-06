package engine

import (
	"errors"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

type Queue struct {
	store data.Storer
}

func NewQueue(store data.Storer) (*Queue, error) {
	if store == nil {
		return nil, errors.New("data store is required for engine.Queue")
	}

	return &Queue{
		store,
	}, nil
}

func (rQueue *Queue) PushTelemetry(newTelemetry ...models.Telemetry) error {
	return nil
}
