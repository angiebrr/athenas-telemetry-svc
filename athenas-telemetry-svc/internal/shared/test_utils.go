package shared

import (
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

func ValidTelemetry() models.Telemetry {
	randTemp := rand.Float64() * 120.0 // rand temp from 0. to 120.

	data := models.Telemetry{
		DeviceID:  uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: randTemp},
		},
	}

	return data
}
