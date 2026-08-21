package specifications

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ================================================================================================

type Telemetry struct {
	DeviceID  string          `json:"device_id"`
	Timestamp int64           `json:"timestamp"`
	Metrics   []MetricReading `json:"metrics"`
}

type MetricReading struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type TelemetryIngester interface {
	Ingest(telemetry Telemetry) error
}

// ------------------------------------------------------------------------------------------------

// TODO: Add other edge cases for TelemetryIngesterSpec

func TelemetryIngesterSpec(testCtx testing.TB, ingester TelemetryIngester) {
	fakeData := Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []MetricReading{
			{"temp", 30.1},
		},
	}
	err := ingester.Ingest(fakeData)
	assert.NoError(testCtx, err)
}
