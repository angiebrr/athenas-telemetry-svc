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
			err := telemetry.Validate(testCase.Data)
			if testCase.ExpectedError != "" {
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(subTestCtx, err)
			}
		})
	}
}
