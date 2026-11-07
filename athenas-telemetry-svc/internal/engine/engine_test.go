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

// FIXME: Add docstrings, potentially move to the data package, shared package, or make a testutils
// package

type BlockingDataStore struct {
	memoryStore *data.InMemoryStore
	gateChan    chan struct{}
}

func NewBlockingDataStore() *BlockingDataStore {
	memoryStore := data.NewInMemoryStore()
	return &BlockingDataStore{
		memoryStore,
		make(chan struct{}),
	}
}

// Insert adds new telemetry data to the in-memory store.
func (rStore *BlockingDataStore) Insert(newTelemetry ...models.Telemetry) error {
	<-rStore.gateChan // this blocks until closed / released by ReleaseBlock

	return rStore.memoryStore.Insert(newTelemetry...)
}

func (rStore *BlockingDataStore) GetByDeviceID(deviceID string) ([]models.Telemetry, error) {
	return rStore.memoryStore.GetByDeviceID(deviceID)
}

func (rStore *BlockingDataStore) ReleaseBlock() {
	close(rStore.gateChan) // unblock inserts
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

	testCtx.Run("queued readings are async from data insertions", func(subTestCtx *testing.T) {
		blockingStore := NewBlockingDataStore()
		subTestCtx.Cleanup(func() {
			blockingStore.ReleaseBlock()
		})
		engineQ, err := engine.NewQueue(blockingStore)
		require.NoError(subTestCtx, err)

		numReadings := 5
		sampleData, deviceID := shared.ValidTelemetryN(numReadings)

		// we basically make sure here that PushTelemetry doesn't wait on the data store
		pushErrChan := make(chan error, 1)
		go func() {
			err := engineQ.PushTelemetry(sampleData...)
			pushErrChan <- err
		}()
		select {
		case err := <-pushErrChan:
			require.NoError(subTestCtx, err)
		case <-time.After(1 * time.Second):
			require.Fail(subTestCtx, "PushTelemetry doesn't asynchronously insert data")
		}

		// ensure the telemetry wasn't inserted
		// (since inserting is blocked until the end of the test)
		_, err = blockingStore.GetByDeviceID(deviceID)
		assert.ErrorIs(subTestCtx, err, data.NotFoundError{DeviceID: deviceID})
	})
}
