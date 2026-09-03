// Package ingest owns the domain rules for accepting telemetry data.
package ingest

import (
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Ingest is the internal domain logic for ingesting telemetry data.
func Ingest(data models.Telemetry, dataStore data.TelemetryDataStorer) error {
	if err := Validate(data); err != nil {
		return err
	}

	if err := dataStore.Insert(data); err != nil {
		return err
	}

	return nil
}
