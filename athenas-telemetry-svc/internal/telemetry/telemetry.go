// Package telemetry owns the domain rules for accepting telemetry data.
package telemetry

import "github.com/angiebrr/athenas-telemetry-svc/internal/data"

// ================================================================================================

// Service represents the internal domain logic for ingesting and querying data in the telemetry
// service, and implements the interfaces for the Telemetry spec.
type Service struct {
	dataStore data.Storer // Where the telemetry is stored and queried from
}

// NewService creates a new Service instance with the given data.Storer as the backing data store
// for managing telemetry.
func NewService(dataStore data.Storer) *Service {
	return &Service{dataStore}
}
