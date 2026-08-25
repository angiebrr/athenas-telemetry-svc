package ingest

import "github.com/angiebrr/athenas-telemetry-svc/models"

// ===============================================================================================

func Ingest(data models.Telemetry) error {
	if err := Validate(data); err != nil {
		return err
	}

	// TODO: do something with the telemetry

	return nil
}
