package ingest_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/ingest"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

// IngestAdapter wraps the Ingest domain logic so it can implement the TelemetryIngester interface
// so we can use the specification tests to verify that the Ingest domain logic behaves as expected.
type IngestAdapter struct{}

func (adapter IngestAdapter) Ingest(data models.Telemetry) error {
	return ingest.Ingest(data)
}

// ------------------------------------------------------------------------------------------------

func TestIngest(testCtx *testing.T) {
	specifications.TelemetryIngesterSpec(testCtx, IngestAdapter{})
}
