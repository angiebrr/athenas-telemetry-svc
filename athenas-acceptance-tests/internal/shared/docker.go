// Package shared provides common utilities for acceptance tests, such as starting a Docker container
// for the telemetry service.
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

// buildLogWriter is a simple helper to bridge io.Writer to the testing Log func testcontainers uses
type buildLogWriter struct{ testCtx testing.TB }

// buildLogWriter.Write calls Log from the wrapped testing.TB instance with the given data
func (rWriter buildLogWriter) Write(data []byte) (int, error) {
	rWriter.testCtx.Log(string(data))
	return len(data), nil
}

// containerLogConsumer is a simple helper to bridge testcontainers.LogConsumer to the testing Log func
//
// TODO: make sure logs don't get logged after test exit
type containerLogConsumer struct{ testCtx testing.TB }

// containerLogConsumer.Accept calls Log from the wrapped testing.TB instance with the given log
func (rConsumer containerLogConsumer) Accept(tcLog testcontainers.Log) {
	rConsumer.testCtx.Logf("CONTAINER LOG [%s]: %s", tcLog.LogType, tcLog.Content)
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
// The docker *build* logs will be tied to testCtx's Log func, and the container will be terminated
// on test cleanup.
//
// NOTE: Only build output is bridged. The container's own stdout is not captured, so a service that
// fails after the image builds surfaces only as a wait-strategy timeout with no explanation.
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
			BuildLogWriter: buildLogWriter{testCtx},
		}),
		testcontainers.WithExposedPorts(port),
		testcontainers.WithWaitStrategy(wait.ForListeningPort(port).WithStartupTimeout(5*time.Second)),
		testcontainers.WithLogConsumers(containerLogConsumer{testCtx}),
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
