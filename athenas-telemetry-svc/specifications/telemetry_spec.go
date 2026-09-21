// Package specifications contains the specifications for the telemetry service, which are used to
// verify that the service behaves as expected. These specifications are intended to be used in both
// acceptance tests and unit tests, and are designed to be implemented by a Driver or Adapter.
package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TODO: docstrings for the interfaces here are weak, just restates the name which isn't helpful

// TelemetryIngester is a system interface for describing the ingestion of telemetry data.
type TelemetryIngester interface {
	Ingest(telemetry models.Telemetry) error
}

// TelemetryQuerier is a system interface for describing the querying of telemetry data.
type TelemetryQuerier interface {
	Query(deviceID string) ([]models.Telemetry, error)
}

// ------------------------------------------------------------------------------------------------

// TelemetrySpec creates a contract for how telemetry ingesting and querying should behave for the
// service as whole, and thus should describe the high-level features of the service.
//
// It is intended to be used to:
//
//   - facilitate acceptance tests via an HTTP client + launched container (
//     found in athenas-acceptance-tests module) OR
//   - via "subcuteanous" unit tests that directly use the internal domain logic in this module
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
