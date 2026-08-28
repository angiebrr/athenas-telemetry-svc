package data_test

import (
	"testing"
	"time"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
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
		sampleData := models.Telemetry{
			DeviceID:  "device123",
			Timestamp: time.Now().Unix(),
			Metrics: []models.MetricReading{
				{Name: "temperature", Value: 25.5},
				{Name: "humidity", Value: 60.0},
			},
		}

		// Insert the sample data into the store
		err := store.Insert(sampleData)
		require.NoError(subTestCtx, err)

		// Retrieve the data by device ID
		retrievedData, err := store.GetByDeviceID("device123")
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
}
