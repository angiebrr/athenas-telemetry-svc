package shared

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ================================================================================================

// Simple helper to bridge io.Writer to testing Logf that testcontainers will use
type testWriter struct{ testCtx testing.TB }

func (rWriter testWriter) Write(data []byte) (int, error) {
	rWriter.testCtx.Log(string(data))
	return len(data), nil
}

// ------------------------------------------------------------------------------------------------

func StartDockerServer(
	testCtx testing.TB,
	port string,
	urlProto string,
) (containerEndpoint string) {
	ctx := context.Background()

	// set up docker container for acceptance test using the server dockerfile
	container, err := testcontainers.Run(
		ctx,
		"", // Empty string for image name since we're building from Dockerfile
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{
			Context:        filepath.Join("..", "..", "..", "..", "src"),
			Dockerfile:     filepath.Join("..", "deploy", "Dockerfile"),
			BuildLogWriter: testWriter{testCtx},
		}),
		testcontainers.WithExposedPorts(port),
		testcontainers.WithWaitStrategy(wait.ForListeningPort(port).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(testCtx, err)

	// schedule cleanup to terminate container after test runs
	testCtx.Cleanup(func() {
		require.NoError(testCtx, testcontainers.TerminateContainer(container))
	})

	// resolve host + mapped port from container (may be dynamically assigned)
	containerEndpoint, err = container.PortEndpoint(ctx, port, urlProto)
	require.NoError(testCtx, err)

	testCtx.Logf("STARTED! Container endpoint: %s", containerEndpoint)

	return containerEndpoint
}
