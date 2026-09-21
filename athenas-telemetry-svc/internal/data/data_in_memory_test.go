package data_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/data"
)

// ================================================================================================

// TestInMemoryStore creates an in-memory data store and runs the data store contract tests against it
func TestInMemoryStore(testCtx *testing.T) {
	store := data.NewInMemoryStore()
	require.NotNil(testCtx, store, "InMemoryStore should not be nil")

	StorerContract(testCtx, store)
}
