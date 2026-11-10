// Package engine handles managing the ingestion of telemetry data at scale
package engine

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TODO: Pass in logger instance to struct here and in the rest of the service

// FIXME: These numbers are very temporary and should be tuned
const (
	DefaultQueueCapacity = 50
	DefaultQueueWorkers  = 10
)

var (
	ErrQueueFull = errors.New("engine queue is full and cannot accept any new work, try again later")
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
func WithWorkers(workers int) OptionFn {
	return func(opts *Options) {
		opts.Workers = workers
	}
}

// ------------------------------------------------------------------------------------------------

// Queue represents a data structure that allows callers to insert pending telemetry asynchronously
type Queue struct {
	store        data.Storer
	readingsChan chan models.Telemetry
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

	engineQ := &Queue{
		store,
		make(chan models.Telemetry, opts.Capacity),
	}
	engineQ.startWork(opts.Workers)
	slog.Info("queue started", "num workers", opts.Workers, "queue capacity", opts.Capacity)

	return engineQ, nil
}

// ------------------------------------------------------------------------------------------------

// Queue.PushTelemetry
// TODO: Make docstring, it's not worth doing until it's fully implemented
func (rQueue *Queue) PushTelemetry(readings ...models.Telemetry) error {
	for _, reading := range readings {
		if err := rQueue.pushWork(reading); err != nil {
			return err
		}
	}

	return nil
}

func (rQueue *Queue) startWork(numWorkers int) {
	for range numWorkers {
		go rQueue.consumeWork()
	}
}

func (rQueue *Queue) consumeWork() {
	for reading := range rQueue.readingsChan {
		if err := rQueue.store.Insert(reading); err != nil {
			// TODO: collect errors somewhere maybe like sentry?
			slog.Error("failed to insert to data store", "err", err)
		}
	}
}

func (rQueue *Queue) pushWork(reading models.Telemetry) error {
	select {
	case rQueue.readingsChan <- reading:
		slog.Debug("reading pushed!")
		return nil
	default:
		slog.Debug("queue full, rejecting reading")
		return ErrQueueFull
	}
}

// FIXME: Add shutdown, context, etc.
