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

// TelemetryIngester is a system interface that is intended to be implemented for use in acceptance
// tests (via a Driver) or in unit tests (via an Adapter)
type TelemetryIngester interface {
	Ingest(telemetry models.Telemetry) error
}

// TelemetryQuerier is a system interface that is intended to be implemented for use in acceptance
// tests (via a Driver) or in unit tests (via an Adapter)
type TelemetryQuerier interface {
	Query(deviceID string) ([]models.Telemetry, error)
}

// ------------------------------------------------------------------------------------------------

// TelemetrySpec verifies that telemetry ingesting and querying works as expected.
func TelemetrySpec(testCtx *testing.T, ingester TelemetryIngester, querier TelemetryQuerier) {
	// TODO: This asserts that invalid input fails, but not *how* it fails. Ingest has exactly one
	// failure mode today, so "any error" and "validation error" describe the same set of outcomes.
	// The moment M2's dispatch ring adds a second class, this assertion starts hiding real bugs --
	// that is the trigger to give the drivers error classification, not a date on a calendar.

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
