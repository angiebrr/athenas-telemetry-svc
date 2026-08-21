package main

import (
	"testing"

	"athenas-telemetry-acceptance-tests/internal/shared"
)

func TestGreeterServer(testCtx *testing.T) {
	// This acceptance test takes a long time, so skip if short running tests are desired
	if testing.Short() {
		testCtx.Skip()
	}

	_ = shared.StartDockerServer(
		testCtx,
		"8080/tcp",
		"http",
	)
}
