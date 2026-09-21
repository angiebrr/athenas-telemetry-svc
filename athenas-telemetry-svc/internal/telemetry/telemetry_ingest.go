package telemetry

import (
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// Service.Ingest is the internal domain logic for ingesting telemetry data.
func (rSvc *Service) Ingest(dataToIngest models.Telemetry) error {
	if err := Validate(dataToIngest); err != nil {
		return err
	}

	return rSvc.dataStore.Insert(dataToIngest)
}
