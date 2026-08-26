// Package ingest owns the domain rules for accepting telemetry data.
package ingest

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ================================================================================================

// Ingest is the internal domain logic for ingesting telemetry data.
func Ingest(data models.Telemetry) error {
	if err := Validate(data); err != nil {
		return err
	}

	// TODO: do something with the telemetry

	return nil
}
