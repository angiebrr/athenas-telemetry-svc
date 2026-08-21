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

func TelemetryIngesterSpec(testCtx *testing.T, ingester TelemetryIngester) {
	fakeData := Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []MetricReading{
			{"temp", 30.1},
		},
	}

	cases := []struct {
		Name          string
		Data          Telemetry
		ExpectedError string
	}{
		{
			Name: "simple metric reading",
			Data: fakeData,
		},
		{
			Name: "invalid metric reading: no metrics",
			Data: func(input Telemetry) Telemetry {
				input.Metrics = nil
				return input
			}(fakeData),
			ExpectedError: "missing metrics",
		},
		{
			Name: "invalid metric reading: no device ID",
			Data: func(input Telemetry) Telemetry {
				input.DeviceID = ""
				return input
			}(fakeData),
			ExpectedError: "missing device ID",
		},
		{
			Name: "invalid metric reading: invalid timestamp",
			Data: func(input Telemetry) Telemetry {
				input.Timestamp = -1
				return input
			}(fakeData),
			ExpectedError: "invalid timestamp",
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
