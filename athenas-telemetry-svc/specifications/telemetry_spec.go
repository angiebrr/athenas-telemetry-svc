// Package specifications contains the specifications for the telemetry service, which are used to
// verify that the service behaves as expected.
//
// These specifications and interfaces are expected to be used by things like an HTTP driver
// that reaches out to the service or an internal adapter that handles the main domain logic.
package specifications

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TelemetryIngester is a port driven by the TelemetrySpec that collects telemetry data.
//
// An error will be returned if invalid telemetry is given or if the data can't be queued and
// should be retried, and nil will be returned if it was able to ingest the data and queue it for
// processing so it will become queryable.
//
// NOTE: Since we're queueing it for processing, a nil return means that querying the data
// immediately after ingestion may not return the data AND that the data could potentially be lost,
// but in most cases it should be queryable after it's processed async.
type TelemetryIngester interface {
	Ingest(telemetry models.Telemetry) error
}

// TelemetryQuerier is a port driven by the TelemetrySpec that retrieves telemetry data that
// has been collected by a TelemetryIngester in the past.
//
// All telemetry for the given device ID will be returned, but it will return an error if no
// data was found.
type TelemetryQuerier interface {
	Query(deviceID string) ([]models.Telemetry, error)
}

// ------------------------------------------------------------------------------------------------

// TelemetrySpec creates a contract for how telemetry ingesting and querying should behave for the
// service as whole, and thus should describe the high-level features of the service.
//
// It is used in a variety of ways to validate this contract, like running acceptance tests against
// the service or running internal unit tests against the internal domain logic.
func TelemetrySpec(testCtx *testing.T, ingester TelemetryIngester, querier TelemetryQuerier) {
	testCtx.Run("valid telemetry is successfully ingested and eventually stored", func(subTestCtx *testing.T) {
		fakeData := shared.ValidTelemetry()

		err := ingester.Ingest(fakeData)
		assert.NoError(subTestCtx, err)

		expectedResultsLen := 1
		queryUntilExpectedLen(subTestCtx, fakeData.DeviceID, querier, expectedResultsLen)
	})

	testCtx.Run("invalid telemetry is rejected on ingest", func(subTestCtx *testing.T) {
		fakeData := shared.ValidTelemetry()
		fakeData.DeviceID = "" // invalid telemetry: no device ID

		err := ingester.Ingest(fakeData)
		assert.ErrorContains(subTestCtx, err, "missing device ID")
	})

	testCtx.Run("querying non-existent telemetry results in error", func(subTestCtx *testing.T) {
		fakeData := shared.ValidTelemetry()
		results, err := querier.Query(fakeData.DeviceID)
		assert.ErrorContains(subTestCtx, err, "data not found")
		assert.Equal(subTestCtx, 0, len(results))
	})
}

// queryUntilExpectedLen is a helper func that calls querier.Query() every 10ms for 5s until the
// querier returns the expected number of results with the target device ID.
func queryUntilExpectedLen(
	testCtx *testing.T,
	targetDeviceID string,
	querier TelemetryQuerier,
	expectedResultsLen int,
) {
	testCtx.Helper()
	queryWaitFor := time.Second * 5
	queryTick := time.Millisecond * 10

	assert.EventuallyWithT(testCtx, func(collectCtx *assert.CollectT) {
		queryResults, queryErr := querier.Query(targetDeviceID)

		require.NoError(collectCtx, queryErr)
		require.Equal(collectCtx, expectedResultsLen, len(queryResults))
		assert.Equal(collectCtx, targetDeviceID, queryResults[0].DeviceID)
	}, queryWaitFor, queryTick)
}
