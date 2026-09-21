// Package telemetry owns the domain rules for accepting telemetry data.
package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Ingest is the internal domain logic for ingesting telemetry data.
func Ingest(dataToIngest models.Telemetry, dataStore data.Storer) error {
	if err := Validate(dataToIngest); err != nil {
		return err
	}

	return dataStore.Insert(dataToIngest)
}
