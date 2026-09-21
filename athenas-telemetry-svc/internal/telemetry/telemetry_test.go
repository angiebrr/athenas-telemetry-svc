package telemetry_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

// IngestAdapter wraps the Ingest domain logic so it can implement the TelemetryIngester interface
// so we can use the specification tests to verify that the Ingest domain logic behaves as expected.
type IngestAdapter struct {
	dataStore data.Storer
}

// IngestAdapter.Ingest satisfies TelemetryIngester by delegating to the package-level Ingest.
func (rAdapter IngestAdapter) Ingest(newData models.Telemetry) error {
	return telemetry.Ingest(newData, rAdapter.dataStore)
}

// QueryAdapter wraps the Query domain logic so it can implement the TelemetryQuerier interface
// so we can use the specification tests to verify that the Query domain logic behaves as expected.
type QueryAdapter struct {
	dataStore data.Storer
}

// QueryAdapter.Query satisfies TelemetryQuerier by delegating to the package-level Query.
func (rAdapter QueryAdapter) Query(deviceID string) ([]models.Telemetry, error) {
	return telemetry.Query(deviceID, rAdapter.dataStore)
}

// ------------------------------------------------------------------------------------------------

// TestTelemetryActions runs the TelemetrySpec against the Ingest and Query domain logic as unit
// teststo verify that the system behaves as expected.
func TestTelemetryActions(testCtx *testing.T) {
	dataStore := data.NewInMemoryStore()
	specifications.TelemetrySpec(
		testCtx,
		IngestAdapter{dataStore},
		QueryAdapter{dataStore},
	)
}
