// Package shared provides common utilities for acceptance tests, such as starting a Docker container
// for the telemetry service.
package shared

import (
	"bytes"
	"context"
	"path/filepath"
	"runtime"
	"sync"
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

// containerLogConsumer bridges testcontainers.LogConsumer to the testing Log func.
//
// NOTE: Accept runs on a testcontainers goroutine that can outlive the test, and testing.TB.Log
// panics once its test completes. The lock is held across the Log call, not merely checked before
// it, and checking first leaves a window for the test to finish in between.
type containerLogConsumer struct {
	testCtx testing.TB

	mutex sync.Mutex
	done  bool
}

// containerLogConsumer.Accept logs one container line, unless the test has already finished.
func (rConsumer *containerLogConsumer) Accept(tcLog testcontainers.Log) {
	rConsumer.mutex.Lock()
	defer rConsumer.mutex.Unlock()

	if rConsumer.done {
		return
	}

	// lines arrive with their trailing newline, which Logf would double up
	rConsumer.testCtx.Logf(
		"CONTAINER LOG [%s]: %s",
		tcLog.LogType,
		bytes.TrimRight(tcLog.Content, "\n"),
	)
}

// containerLogConsumer.stop makes every subsequent Accept a no-op.
func (rConsumer *containerLogConsumer) stop() {
	rConsumer.mutex.Lock()
	defer rConsumer.mutex.Unlock()

	rConsumer.done = true
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
// Both the docker build logs and the running container's stdout/stderr are tied to testCtx's Log
// func, and the container will be terminated on test cleanup.
//
// NOTE: These are separate paths -- BuildLogWriter covers the image build, a LogConsumer covers the
// running container. A service that builds fine but dies at startup only appears in the second.
func StartDockerServer(
	testCtx testing.TB,
	port string,
	urlProto string,
) (containerEndpoint string) {
	ctx := context.Background()

	// Bridge the container's own stdout/stderr into the test log.
	//
	// NOTE: Registered before the container's cleanup so it runs after it since cleanup is
	// last-added-first-called, and terminating the container drains buffered shutdown lines that
	// closing the gate early would discard.
	logConsumer := &containerLogConsumer{testCtx: testCtx}
	testCtx.Cleanup(logConsumer.stop)

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
		testcontainers.WithLogConsumers(logConsumer),
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
