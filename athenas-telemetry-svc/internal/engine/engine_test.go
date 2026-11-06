package engine_test

import (
	"testing"
	"time"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/engine"
	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================================================

func TestEngineQueue(testCtx *testing.T) {
	// set up data store and engine
	//
	// TODO: Will there be other kinds of engine queues? Should this be a contract?
	store := data.NewInMemoryStore()
	engineQ, err := engine.NewQueue(store)
	require.NoError(testCtx, err)

	testCtx.Run("queued readings eventually reach the data store", func(subTestCtx *testing.T) {
		numReadings := 5
		sampleData, deviceID := shared.ValidTelemetryN(numReadings)

		err := engineQ.PushTelemetry(sampleData...)
		require.NoError(subTestCtx, err)

		queryStoreUntilExpectedLen(subTestCtx, deviceID, store, numReadings)
	})

}

func queryStoreUntilExpectedLen(
	testCtx *testing.T,
	targetDeviceID string,
	store data.Storer,
	expectedResultsLen int,
) {
	testCtx.Helper()
	dataWaitFor := time.Second * 5
	dataTick := time.Millisecond * 10

	assert.EventuallyWithT(testCtx, func(collectCtx *assert.CollectT) {
		queryResults, queryErr := store.GetByDeviceID(targetDeviceID)

		require.NoError(collectCtx, queryErr)
		require.Equal(collectCtx, expectedResultsLen, len(queryResults))
		assert.Equal(collectCtx, targetDeviceID, queryResults[0].DeviceID)
	}, dataWaitFor, dataTick)
}
