// Package specifications contains the specifications for the telemetry service, which are used to
// verify that the service behaves as expected. These specifications are intended to be used in both
// acceptance tests and unit tests, and are designed to be implemented by a Driver or Adapter.
package specifications

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TelemetryIngester is a system interface that is intended to be implemented for use in acceptance
// tests (via a Driver) or in unit tests (via an Adapter)
type TelemetryIngester interface {
	Ingest(telemetry models.Telemetry) error
}

// ------------------------------------------------------------------------------------------------

// TelemetryIngesterSpec verifies that an Ingester behaves as expected
func TelemetryIngesterSpec(testCtx *testing.T, ingester TelemetryIngester) {
	fakeData := models.Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: 30.1},
		},
	}

	// TODO: This doesn't seem like enough?

	testCtx.Run("ingest valid telemetry", func(subTestCtx *testing.T) {
		err := ingester.Ingest(fakeData)
		assert.NoError(subTestCtx, err)
	})

	testCtx.Run("ingest invalid telemetry", func(subTestCtx *testing.T) {
		err := ingester.Ingest(models.Telemetry{})
		assert.Error(subTestCtx, err)
	})
}
