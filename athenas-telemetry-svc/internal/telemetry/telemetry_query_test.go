package telemetry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
)

// ===============================================================================================

// TestQueryDeviceIDs unit tests the Query function using various device IDs, as it's out of scope
// for the specification tests to cover this function in detail
//
// (see telemetry_test.go that uses the specs as unit tests)
func TestQueryDeviceIDs(testCtx *testing.T) {
	testCases := []struct {
		Name          string
		DeviceID      string
		SetupFn       func(dataStore data.Storer) error
		ExpectedError string
	}{
		{
			Name:          "empty device ID",
			DeviceID:      "",
			ExpectedError: "missing device ID",
		},
		{
			Name:     "valid device ID",
			DeviceID: "device-123",
			SetupFn: func(dataStore data.Storer) error {
				validData := shared.ValidTelemetry()
				validData.DeviceID = "device-123"
				return dataStore.Insert(validData)
			},
			ExpectedError: "",
		},
	}

	dataStore := data.NewInMemoryStore()

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			// setup data for test case
			if testCase.SetupFn != nil {
				err := testCase.SetupFn(dataStore)
				require.NoError(subTestCtx, err, "failed to set up test case")
			}

			// execute the query and validate our expectations for the test
			results, err := telemetry.Query(testCase.DeviceID, dataStore)
			if testCase.ExpectedError != "" {
				assert.Error(subTestCtx, err)
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(subTestCtx, err)
				// the results don't matter for this test, since that's covered in spec tests, but
				// we can at least assert that we got some results back for a valid device ID
				assert.NotNil(subTestCtx, results)
			}
		})
	}
}
