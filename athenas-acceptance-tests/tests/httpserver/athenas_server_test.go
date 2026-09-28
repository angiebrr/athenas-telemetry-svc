package httpserver_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/angiebrr/athenas-telemetry-svc/specifications"
	"github.com/stretchr/testify/require"

	"athenas-telemetry-acceptance-tests/internal/drivers/httpserver"
	"athenas-telemetry-acceptance-tests/internal/shared"
)

// ================================================================================================

func TestAthenasTelemetryServer(testCtx *testing.T) {
	// This acceptance test takes a long time, so skip if short running tests are desired
	if testing.Short() {
		testCtx.Skip()
	}

	// Run the telemetry service locally via a docker image
	containerEndpoint := shared.StartDockerServer(
		testCtx,
		"8080/tcp",
		"http",
	)

	// Make a test driver to run HTTP reqs against the telemetry service
	driver, err := httpserver.NewDriver(containerEndpoint, &http.Client{
		Timeout: 1 * time.Second,
	})
	require.NoError(testCtx, err)

	// Run the driver against the telemetry spec tests
	specifications.TelemetrySpec(testCtx, driver, driver)
}
