package main

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/specifications"

	"athenas-telemetry-acceptance-tests/internal/shared"
)

// ================================================================================================

// TODO: make a real HTTP driver here
type DummyDriver struct {
}

func (rDriver *DummyDriver) Ingest(telemetry specifications.Telemetry) error {
	return nil
}

func TestAthenasTelemetryServer(testCtx *testing.T) {
	// This acceptance test takes a long time, so skip if short running tests are desired
	if testing.Short() {
		testCtx.Skip()
	}

	// Run the telemetry service locally via a docker image
	_ = shared.StartDockerServer(
		testCtx,
		"8080/tcp",
		"http",
	)

	// Make a test driver to run HTTP reqs against the telemetry service
	driver := &DummyDriver{}

	// Run the driver against the telemetry spec tests
	specifications.TelemetryIngesterSpec(testCtx, driver)
}
