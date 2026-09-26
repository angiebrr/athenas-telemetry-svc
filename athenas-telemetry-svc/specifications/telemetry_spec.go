// Package specifications contains the specifications for the telemetry service, which are used to
// verify that the service behaves as expected.
//
// These specifications and interfaces are expected to be used by things like an HTTP driver
// that reaches out to the service or an internal adapter that handles the main domain logic.
package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TelemetryIngester is a port driven by the TelemetrySpec that collects telemetry data.
//
// An error will be returned if invalid telemetry is given, and nil will be returned if it was able
// to ingest the data.
//
// TODO: It's a little vague if "able to ingest" means if it can be queried immediately after or
// not, but this will be firmed up after future work which uses ring buffers to process lots
// of data at once.
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
	testCtx.Run("valid telemetry is successfully ingested", func(subTestCtx *testing.T) {
		fakeData := shared.ValidTelemetry()

		err := ingester.Ingest(fakeData)
		assert.NoError(subTestCtx, err)

		results, err := querier.Query(fakeData.DeviceID)
		assert.NoError(subTestCtx, err)
		require.Equal(subTestCtx, 1, len(results))
		assert.Equal(subTestCtx, fakeData.DeviceID, results[0].DeviceID)
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
