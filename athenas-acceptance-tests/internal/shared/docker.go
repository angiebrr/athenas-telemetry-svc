package shared

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ================================================================================================

// Simple helper to bridge io.Writer to testing Log func that testcontainers will use
type testWriter struct{ testCtx testing.TB }

// Write calls Log from the wrapped testing.TB instance with the given data
func (rWriter testWriter) Write(data []byte) (int, error) {
	rWriter.testCtx.Log(string(data))
	return len(data), nil
}

// ------------------------------------------------------------------------------------------------

// repoRoot returns the path of the main repository, which is useful for getting the dockerfile
// in the deploy/ folder.
//
// NOTE: This is a little fragile since it will break if this source file is moved, but it's easier
// than passing through an env var or walking the fs until we find a Dockerfile
func repoRoot() string {
	// get the path of the dir for this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed for current file")
	}
	currDir := filepath.Dir(filename)

	// get our repo path by moving relative to where this script is located
	repoRoot := filepath.Join(currDir, "..", "..", "..")
	return filepath.Clean(repoRoot)
}

// ------------------------------------------------------------------------------------------------

// StartDockerServer runs the athena telemetry service using the Dockerfile in deploy/ so our
// acceptance test can run against it.
//
// The dockerfile logs will be tied to testCtx's Log func, and the container will be terminated on
// test cleanup.
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
			Context:        repoRoot(),
			Dockerfile:     filepath.Join("deploy", "Dockerfile"),
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
