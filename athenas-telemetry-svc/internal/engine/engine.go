// Package engine handles managing the ingestion of telemetry data at scale
package engine

import (
	"errors"
	"fmt"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// FIXME: These numbers are very temporary and should be tuned
const (
	DefaultQueueCapacity = 50
	DefaultQueueWorkers  = 10
)

var (
	QueueFullError = errors.New("engine queue is full and cannot accept any new work, try again later")
)

// ------------------------------------------------------------------------------------------------

// OptionFn is the helper func for modifying the Options input struct
type OptionFn func(*Options)

// Options defines the args that can be passed into NewQueue that can be modified by OptionFn
type Options struct {
	Capacity int
	Workers  int
}

// Options.Validate ensures that the input values are correct.
//
// Although we set defaults for optional params, we still validate everything since they
// could be overridden with input OptionFns
func (rOpts *Options) Validate() error {
	if rOpts.Capacity <= 0 {
		return errors.New("capacity must be greater than 0")
	}
	if rOpts.Workers <= 0 {
		return errors.New("workers must be greater than 0")
	}

	return nil
}

// WithCapacity sets the capacity of the queue before it starts rejecting work being pushed onto it.
// By default, it's set to engine.DefaultQueueCapacity
func WithCapacity(capacity int) OptionFn {
	return func(opts *Options) {
		opts.Capacity = capacity
	}
}

// WithWorkers sets the number of workers / go-routines that will be created to process the queue.
// By default, it's set to engine.DefaultQueueWorkers
func WithWorkers(capacity int) OptionFn {
	return func(opts *Options) {
		opts.Capacity = capacity
	}
}

// ------------------------------------------------------------------------------------------------

// Queue represents a data structure that allows callers to insert pending telemetry asynchronously
type Queue struct {
	store data.Storer
}

// NewQueue creates a new Queue instance with the given data.Storer as the backing data store
// that the queue workers will use to manage telemetry.
//
// It accepts a slice of OptionFn funcs that can be used to override the default values.
func NewQueue(store data.Storer, inOptFns ...OptionFn) (*Queue, error) {
	opts := &Options{
		Capacity: DefaultQueueCapacity,
		Workers:  DefaultQueueWorkers,
	}
	for _, optFn := range inOptFns {
		optFn(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options to engine.NewQueue: %w", err)
	}

	if store == nil {
		return nil, errors.New("data store is required for engine.Queue")
	}

	// FIXME: do something with the capacity and num workers

	return &Queue{
		store,
	}, nil
}

// ------------------------------------------------------------------------------------------------

// Queue.PushTelemetry
// TODO: Make docstring, it's not worth doing until it's fully implemented
func (rQueue *Queue) PushTelemetry(newTelemetry ...models.Telemetry) error {
	// FIXME: This needs to actually use workers and a queue, but need a better test for that first
	go func() {
		_ = rQueue.store.Insert(newTelemetry...)
	}()
	return nil
}
