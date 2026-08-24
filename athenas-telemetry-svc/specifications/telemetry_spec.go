package specifications

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// TelemetryIngester is a system interface that is intended to be implemented for use in acceptance
// tests (via a Driver) or in unit tests (via an Adapter)
type TelemetryIngester interface {
	Ingest(telemetry models.Telemetry) error
}

// ------------------------------------------------------------------------------------------------

// TelemetryIngesterSpec verifies that an Ingester behaves as expected
func TelemetryIngesterSpec(testCtx *testing.T, ingester TelemetryIngester) {
	fakeData := models.Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: 30.1},
		},
	}

	// TODO: these errors are specific to gin/JSON, but getting the validator to use custom english responses is A WHOLE THING
	// so we're doing that later
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
			ExpectedError: "Metrics",
		},
		{
			Name: "invalid metric reading: no device ID",
			Data: func(input models.Telemetry) models.Telemetry {
				input.DeviceID = ""
				return input
			}(fakeData),
			ExpectedError: "DeviceID",
		},
		{
			Name: "invalid metric reading: invalid timestamp",
			Data: func(input models.Telemetry) models.Telemetry {
				input.Timestamp = -1
				return input
			}(fakeData),
			ExpectedError: "Timestamp",
		},
	}

	for _, testCase := range cases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			err := ingester.Ingest(testCase.Data)
			if testCase.ExpectedError != "" {
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(testCtx, err)
			}
		})
	}
}
