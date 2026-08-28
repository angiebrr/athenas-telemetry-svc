package ingest_test

import (
	"testing"
	"time"

	"github.com/angiebrr/athenas-telemetry-svc/internal/ingest"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/angiebrr/athenas-telemetry-svc/specifications"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

// IngestAdapter wraps the Ingest domain logic so it can implement the TelemetryIngester interface
// so we can use the specification tests to verify that the Ingest domain logic behaves as expected.
type IngestAdapter struct{}

// IngestAdapter.Ingest satisfies TelemetryIngester by delegating to the package-level Ingest.
func (IngestAdapter) Ingest(data models.Telemetry) error {
	return ingest.Ingest(data)
}

type QueryAdapter struct{}

func (QueryAdapter) Query(deviceID string) ([]models.Telemetry, error) {
	return ingest.Query(deviceID)
}

func TestIngest(testCtx *testing.T) {
	specifications.TelemetrySpec(testCtx, IngestAdapter{}, QueryAdapter{})
}

// ------------------------------------------------------------------------------------------------

func TestIngest_Validate(testCtx *testing.T) {
	// TODO: Make a shared "valid telemetry" function so we can use it in both the spec tests and these unit tests
	fakeData := models.Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: 30.1},
		},
	}

	cases := []struct {
		Name          string
		Data          models.Telemetry
		ExpectedError string
	}{
		{
			Name: "simple metric reading",
			Data: fakeData,
		},
		{
			Name: "invalid metric reading: no metrics",
			Data: func(input models.Telemetry) models.Telemetry {
				input.Metrics = nil
				return input
			}(fakeData),
			ExpectedError: "missing metrics",
		},
		{
			Name: "invalid metric reading: no device ID",
			Data: func(input models.Telemetry) models.Telemetry {
				input.DeviceID = ""
				return input
			}(fakeData),
			ExpectedError: "missing device ID",
		},
		{
			Name: "invalid metric reading: invalid timestamp",
			Data: func(input models.Telemetry) models.Telemetry {
				input.Timestamp = -1
				return input
			}(fakeData),
			ExpectedError: "invalid timestamp",
		},
		{
			Name: "invalid metric reading: missing metric name",
			Data: func(input models.Telemetry) models.Telemetry {
				input.Metrics = []models.MetricReading{
					{Name: "", Value: 30.1},
				}
				return input
			}(fakeData),
			ExpectedError: "missing metric name",
		},
	}

	for _, testCase := range cases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			err := ingest.Validate(testCase.Data)
			if testCase.ExpectedError != "" {
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(subTestCtx, err)
			}
		})
	}
}
