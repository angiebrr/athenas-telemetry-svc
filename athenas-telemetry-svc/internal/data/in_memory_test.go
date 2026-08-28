package data_test

import (
	"sync"
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================================================

func TestInMemoryDataStore(testCtx *testing.T) {
	// Create a new in-memory data store for testing
	store := data.NewInMemoryDataStore()
	require.NotNil(testCtx, store, "InMemoryDataStore should not be nil")

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
		err := store.Insert(sampleData1)
		require.NoError(subTestCtx, err)
		err = store.Insert(sampleData2)
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

				// insert and read data for new device (contend on map)
				newData := shared.ValidTelemetry()
				insertGetAndIterate(subTestCtx, newData, store)

				// insert and read same device over and over (contend on map + same device array)
				insertGetAndIterate(subTestCtx, contendedData, store)
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
}

func insertGetAndIterate(testCtx *testing.T, newData models.Telemetry, store data.TelemetryDataStorer) {
	testCtx.Helper()

	// Insert the data and make sure it exists
	// Note that we use "assert" instead of "require" so the go-routine doesn't silently die
	err := store.Insert(newData)
	if ok := assert.NoError(testCtx, err); !ok {
		return
	}
	results, err := store.GetByDeviceID(newData.DeviceID)
	if ok := assert.NoError(testCtx, err); !ok {
		return
	}

	// Read data from slice to test thread-safe reads and make sure we aren't leaking backing data
	// that isn't thread-safe
	for _, currData := range results {
		assert.Equal(testCtx, newData.DeviceID, currData.DeviceID)
	}
}
