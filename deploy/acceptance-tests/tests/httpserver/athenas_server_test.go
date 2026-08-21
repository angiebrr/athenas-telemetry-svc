package main

import (
	"testing"

	"athenas-telemetry-acceptance-tests/internal/shared"
)

// ================================================================================================

func TestAthenasTelemetryServer(testCtx *testing.T) {
	// This acceptance test takes a long time, so skip if short running tests are desired
	if testing.Short() {
		testCtx.Skip()
	}

	_ = shared.StartDockerServer(
		testCtx,
		"8080/tcp",
		"http",
	)

	// TODO: add driver for HTTP that specification tests use

	// TODO: add specification test
}
