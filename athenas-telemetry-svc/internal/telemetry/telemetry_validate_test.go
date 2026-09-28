package telemetry_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
	"github.com/angiebrr/athenas-telemetry-svc/models"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

// TestTelemetry_Validate tests the Validate function used to ensure that telemetry data is valid
// before it is ingested into the system.
//
// Note that this is separate from spec tests since it should be tested in isolation and is an
// implementation detail of the system.
func TestTelemetry_Validate(testCtx *testing.T) {
	fakeData := shared.ValidTelemetry()

	cases := []struct {
		Name          string
		Data          models.Telemetry
		ExpectedError error
	}{
		{
			Name: "simple metric reading",
			Data: fakeData,
		},
		{
			Name:          "invalid metric reading: nil metrics",
			Data:          telemetryWithMetrics(fakeData, nil),
			ExpectedError: telemetry.ErrMissingMetrics,
		},
		{
			Name:          "invalid metric reading: no metrics",
			Data:          telemetryWithMetrics(fakeData, []models.MetricReading{}),
			ExpectedError: telemetry.ErrMissingMetrics,
		},
		{
			Name:          "invalid metric reading: no device ID",
			Data:          telemetryWithDeviceID(fakeData, ""),
			ExpectedError: telemetry.ErrMissingDeviceID,
		},
		{
			Name:          "invalid metric reading: invalid timestamp",
			Data:          telemetryWithTimestamp(fakeData, -1),
			ExpectedError: telemetry.ErrInvalidTimestamp,
		},
		{
			Name:          "invalid metric reading: missing timestamp",
			Data:          telemetryWithTimestamp(fakeData, 0),
			ExpectedError: telemetry.ErrInvalidTimestamp,
		},
		{
			Name: "invalid metric reading: missing metric name",
			Data: telemetryWithMetrics(fakeData, []models.MetricReading{
				{Name: "uptimeSecs", Value: 1000.},
				{Name: "", Value: 30.1},
			}),
			ExpectedError: telemetry.ErrMissingMetricName,
		},
	}

	for _, testCase := range cases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			err := telemetry.Validate(testCase.Data)
			if testCase.ExpectedError != nil {
				assert.ErrorIs(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(subTestCtx, err)
			}
		})
	}
}
