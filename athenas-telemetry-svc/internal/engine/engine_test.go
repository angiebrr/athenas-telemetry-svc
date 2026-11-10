package engine_test

import (
	"testing"
	"time"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/engine"
	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================================================

// TODO: Look into creating a slog instance and pass in testing obj
// logger := slog.New(slog.NewTextHandler(testCtx.Output(), nil))

// FIXME: Potentially move  BlockingDataStore to the data package, shared package, or make a
// testutils package

// BlockingDataStore is a stub that uses an in-memory data store for operations but blocks on
// inserts until the block is released on shutdown.
//
// It is useful for verifying queue behavior in tests.
type BlockingDataStore struct {
	memoryStore *data.InMemoryStore
	gateChan    chan struct{}
}

// NewBlockingDataStore creates a stub data store that tests can use to verify queue behavior.
func NewBlockingDataStore() *BlockingDataStore {
	memoryStore := data.NewInMemoryStore()
	return &BlockingDataStore{
		memoryStore,
		make(chan struct{}),
	}
}

// BlockingDataStore.Insert adds new telemetry data to the in-memory store by calling
// data.InMemoryStore.Insert.
//
// This will block until BlockingDataStore.Shutdown is called.
func (rStore *BlockingDataStore) Insert(newTelemetry ...models.Telemetry) error {
	<-rStore.gateChan
	return rStore.memoryStore.Insert(newTelemetry...)
}

// BlockingDataStore.GetByDeviceID retrieves inserted telemetry from the in-memory store by calling
// data.InMemoryStore.GetByDeviceByID.
//
// It does not block.
func (rStore *BlockingDataStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	return rStore.memoryStore.GetByDeviceID(deviceID)
}

// BlockingDataStore.Shutdown releases the channel blocking BlockingDataStore.Insert so inserts can
// now go through.
func (rStore *BlockingDataStore) Shutdown() {
	close(rStore.gateChan)
}

// ------------------------------------------------------------------------------------------------

// pushTelemetryWithBlockingStore is a test helper to push readings and verify queue behavior with
// BlockingDataStore by creating some data, pushing that data to the queue, and then checking the
// data store.
//
// You can customize what error is expected by PushTelemetry, but, since the data store is blocking
// until the end of the test, retrieving the data after push should always result in a
// data.NotFoundError
func pushTelemetryWithBlockingStore(
	testCtx *testing.T,
	engineQ *engine.Queue,
	blockingStore *BlockingDataStore,
	numReadings int,
	expectedPushError error,
) {
	testCtx.Helper()

	// push N telemetry readings to the queue, where N is the input numReadings
	sampleData, targetDeviceID := shared.ValidTelemetryN(numReadings)
	pushErrChan := make(chan error, 1)
	go func() {
		err := engineQ.PushTelemetry(sampleData...)
		pushErrChan <- err
	}()

	// since blockingStore blocks inserts, we verify PushTelemetry is async, which is why it's in a
	// go-routine and we're timing out after a second
	select {
	case err := <-pushErrChan:
		if expectedPushError == nil {
			require.NoError(testCtx, err)
		} else {
			require.ErrorIs(testCtx, err, expectedPushError)
		}
	case <-time.After(1 * time.Second):
		require.Fail(testCtx, "PushTelemetry doesn't asynchronously insert data or is blocked by something else")
	}

	// ensure the telemetry wasn't inserted
	// (since inserting is blocked until the end of the test)
	_, err := blockingStore.GetByDeviceID(targetDeviceID)
	require.ErrorIs(testCtx, err, data.NotFoundError{DeviceID: targetDeviceID})
}

// ================================================================================================

func TestEngineQueue(testCtx *testing.T) {
	const eventuallyWaitFor = time.Second * 5
	const eventuallyTick = time.Millisecond * 10

	testCtx.Run("queued readings eventually reach the data store", func(subTestCtx *testing.T) {
		memoryStore := data.NewInMemoryStore()
		engineQ, err := engine.NewQueue(memoryStore)
		require.NoError(subTestCtx, err)

		numReadings := 5
		sampleData, deviceID := shared.ValidTelemetryN(numReadings)

		err = engineQ.PushTelemetry(sampleData...)
		require.NoError(subTestCtx, err)

		// call the data store periodically after pushing, ensuring the readings are eventually
		// in the data store
		assert.EventuallyWithT(subTestCtx, func(collectCtx *assert.CollectT) {
			results, err := memoryStore.GetByDeviceID(deviceID)
			require.NoError(collectCtx, err)
			require.Equal(collectCtx, numReadings, len(results))
			assert.Equal(collectCtx, deviceID, results[0].DeviceID)
		}, eventuallyWaitFor, eventuallyTick)
	})

	testCtx.Run("full queues return a 'full' error immediately", func(subTestCtx *testing.T) {
		blockingStore := NewBlockingDataStore()
		subTestCtx.Cleanup(func() {
			blockingStore.Shutdown()
		})

		const maxQItems = 5
		const oneWorker = 1
		engineQ, err := engine.NewQueue(blockingStore,
			engine.WithCapacity(maxQItems),
			engine.WithWorkers(oneWorker),
		)
		require.NoError(subTestCtx, err)

		// push the max number of telemetry readings to the queue
		pushTelemetryWithBlockingStore(
			subTestCtx,
			engineQ,
			blockingStore,
			maxQItems,
			nil, // PushTelemetry shouldn't error
		)

		// then try to push another set of readings, which should immediately error
		pushTelemetryWithBlockingStore(
			subTestCtx,
			engineQ,
			blockingStore,
			maxQItems,
			engine.ErrQueueFull, // PushTelemetry should report a full queue
		)
	})
}
