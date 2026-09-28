package telemetry_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

// TestTelemetryService runs the TelemetrySpec against the Ingest and Query domain logic as unit
// tests to verify that the system behaves as expected.
func TestTelemetryService(testCtx *testing.T) {
	dataStore := data.NewInMemoryStore()
	telemetrySvc := telemetry.NewService(dataStore)

	specifications.TelemetrySpec(
		testCtx,
		telemetrySvc,
		telemetrySvc,
	)
}
