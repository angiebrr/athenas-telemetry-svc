package data_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// insertAndGet is a helper for the go-routines in our thread-safe tests to call when they are
// doing their work. It will insert newData using the given store, and then try to retrieve
// that inserted data from the store.
func insertAndGet(testCtx *testing.T, newData models.Telemetry, store data.Storer) {
	testCtx.Helper()

	// Insert the data and make sure it exists
	// Note that we use "assert" instead of "require" so the go-routine doesn't silently die
	err := store.Insert(newData)
	if ok := assert.NoError(testCtx, err); !ok {
		return
	}
	_, err = store.GetByDeviceID(newData.DeviceID)
	assert.NoError(testCtx, err)
}

// ------------------------------------------------------------------------------------------------

// StorerContract creates a contract for how a data.Storer implementation should behave within
// the internals of the telemetry service.
//
// It is intended to be used in unit tests to verify that a given implementation of the interface
// behaves as expected (e.g. in-memory store, database-backed store, etc.)
func StorerContract(testCtx *testing.T, store data.Storer) {
	testCtx.Run("insert and get data from store", func(subTestCtx *testing.T) {
		// Create a sample telemetry data
		sampleData := shared.ValidTelemetry()

		// Insert the sample data into the store
		err := store.Insert(sampleData)
		require.NoError(subTestCtx, err)

		// Retrieve the data by device ID
		retrievedData, err := store.GetByDeviceID(sampleData.DeviceID)
		require.NoError(subTestCtx, err)

		// Verify that the retrieved data matches the inserted data
		require.Len(subTestCtx, retrievedData, 1)
		insertedData := retrievedData[0]
		assert.Equal(subTestCtx, sampleData.DeviceID, insertedData.DeviceID)
		assert.Equal(subTestCtx, sampleData.Timestamp, insertedData.Timestamp)
		assert.Equal(subTestCtx, sampleData.Metrics, insertedData.Metrics)
	})

	testCtx.Run("get data for non-existent device ID", func(subTestCtx *testing.T) {
		// Attempt to retrieve data for a device ID that doesn't exist
		retrievedData, err := store.GetByDeviceID("nonexistent_device")
		require.Error(subTestCtx, err) // device not found

		// Verify that the retrieved data is empty
		assert.Empty(subTestCtx, retrievedData)
	})

	testCtx.Run("insert multiple data entries for the same device ID", func(subTestCtx *testing.T) {
		// Create multiple sample telemetry data entries for the same device ID
		sampleData1 := shared.ValidTelemetry()
		sampleData2 := shared.ValidTelemetry()
		sampleData2.DeviceID = sampleData1.DeviceID

		// Insert both sample data entries into the store
		err := store.Insert(sampleData1, sampleData2)
		require.NoError(subTestCtx, err)

		// Retrieve the data by device ID
		retrievedData, err := store.GetByDeviceID(sampleData1.DeviceID)
		require.NoError(subTestCtx, err)

		// Verify that both entries are retrieved
		require.Len(subTestCtx, retrievedData, 2)
		assert.Equal(subTestCtx, sampleData1.DeviceID, retrievedData[0].DeviceID)
		assert.Equal(subTestCtx, sampleData2.DeviceID, retrievedData[1].DeviceID)
	})

	testCtx.Run("make multiple inserts concurrently", func(subTestCtx *testing.T) {
		contendedData := shared.ValidTelemetry()

		// Insert and read the sample data concurrently
		numGoRoutines := 100
		var allOpsWg, gateWg sync.WaitGroup
		gateWg.Add(1)
		for range numGoRoutines {
			allOpsWg.Go(func() {
				gateWg.Wait() // wait for all go-routines to be added

				// insert and read data for new device
				newData := shared.ValidTelemetry()
				insertAndGet(subTestCtx, newData, store)

				// insert and read same device over and over
				insertAndGet(subTestCtx, contendedData, store)
			})
		}
		gateWg.Done()   // signal go-routines to run
		allOpsWg.Wait() // wait until they are done

		// Verify the contended device has the number of expected telemetry datapoints
		receivedData, err := store.GetByDeviceID(contendedData.DeviceID)
		require.NoError(subTestCtx, err)
		assert.Len(
			subTestCtx,
			receivedData,
			numGoRoutines,
			"contended device doesn't have the expected number of datapoints",
		)
	})

	testCtx.Run("result slice isn't leaked", func(subTestCtx *testing.T) {
		// make 3 pieces of telemetry with the same device ID
		sampleResults := shared.ValidTelemetryN(3)
		sampleData1 := sampleResults[0]
		sampleData2 := sampleResults[1]
		sampleData3 := sampleResults[2]
		sampleData2.DeviceID = sampleData1.DeviceID
		sampleData3.DeviceID = sampleData1.DeviceID

		// Insert the sample data and try to get it
		err := store.Insert(sampleData1, sampleData2, sampleData3)
		require.NoError(subTestCtx, err)
		results, err := store.GetByDeviceID(sampleData1.DeviceID)
		require.NoError(subTestCtx, err)
		require.Len(subTestCtx, results, 3, "device should have 3 pieces of telemetry")

		// Mutate the local results array
		results[0].DeviceID = "MUTATED"

		// Get the results again, and verify none of the data points have a mutated deviced ID
		results, err = store.GetByDeviceID(sampleData1.DeviceID)
		require.NoError(subTestCtx, err)
		require.Len(subTestCtx, results, 3, "device should have 3 pieces of telemetry")
		for _, currData := range results {
			got := currData.DeviceID
			want := sampleData1.DeviceID
			assert.Equal(subTestCtx, want, got)
		}
	})
}
