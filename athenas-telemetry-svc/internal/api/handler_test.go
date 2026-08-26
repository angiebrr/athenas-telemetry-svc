package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angiebrr/athenas-telemetry-svc/internal/api"
	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// errorResponse mirrors the JSON error envelope the telemetry handlers write so tests can assert on
// the message itself rather than on a raw body string.
type errorResponse struct {
	Error string `json:"error"`
}

// validTelemetry returns telemetry that passes validation, for tests to mutate into invalid shapes.
func validTelemetry() models.Telemetry {
	return models.Telemetry{
		DeviceID:  "12345",
		Timestamp: time.Now().UnixMilli(),
		Metrics: []models.MetricReading{
			{Name: "temp", Value: 30.1},
		},
	}
}

// ------------------------------------------------------------------------------------------------

// newTestServer builds a server instance for tests to serve requests against.
//
// NOTE: The server holds no request state as of M1, so one instance can be shared by every
// subtest. Once it owns a dispatch ring each case needs its own, so they can't leak state into
// each other.
func newTestServer(testCtx testing.TB) *api.Server {
	// quiet gin debug logs during testing
	gin.SetMode(gin.TestMode)

	server := api.NewServer()
	require.NotNil(testCtx, server, "server should not be nil")

	return server
}

// doRequest sends one request directly into the gin engine and returns the recorded response.
func doRequest(
	server *api.Server,
	method string,
	path string,
	body io.Reader,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	return recorder
}

// postTelemetry serializes the given telemetry and POSTs it to the ingest endpoint.
func postTelemetry(
	testCtx testing.TB,
	server *api.Server,
	data models.Telemetry,
) *httptest.ResponseRecorder {
	payload, err := json.Marshal(data)
	require.NoError(testCtx, err, "test telemetry should serialize")

	return doRequest(server, http.MethodPost, api.TelemetryPath, bytes.NewReader(payload))
}

// decodeError pulls the message out of the handler's JSON error envelope.
func decodeError(testCtx testing.TB, recorder *httptest.ResponseRecorder) string {
	var response errorResponse

	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(testCtx, err, "error responses should carry a JSON error envelope")

	return response.Error
}

// ================================================================================================

// TestHandleIngestTelemetry verifies the transport-level behaviour of the ingest endpoint: the
// status codes it maps onto, and the routing it does and does not accept.
//
// These are deliberately the assertions the telemetry spec cannot make. The spec speaks
// models.Telemetry rather than bytes, so it can neither send a malformed body nor tell 202 from 200
// nor 400 from 500.
func TestHandleIngestTelemetry(testCtx *testing.T) {
	server := newTestServer(testCtx)

	// ---

	testCtx.Run("accepts valid telemetry with 202", func(subTestCtx *testing.T) {
		recorder := postTelemetry(subTestCtx, server, validTelemetry())

		assert.Equal(subTestCtx, http.StatusAccepted, recorder.Code)
		assert.Empty(subTestCtx, recorder.Body.String(), "an accepted ingest should return no body")
	})

	// This is the only case in the repo that exercises the ShouldBindJSON failure branch, since the
	// spec's interface has no way to send bytes that aren't valid telemetry.
	testCtx.Run("rejects a malformed body before reaching the domain", func(subTestCtx *testing.T) {
		recorder := doRequest(
			server,
			http.MethodPost,
			api.TelemetryPath,
			bytes.NewBufferString(`{"device_id":`),
		)

		assert.Equal(subTestCtx, http.StatusBadRequest, recorder.Code)
		assert.NotEmpty(
			subTestCtx,
			decodeError(subTestCtx, recorder),
			"a malformed body should still report an error",
		)
	})

	// The concrete, HTTP-specific half of the error classification the spec asserts only
	// loosely (as "some error") while Ingest still has a single failure mode.
	testCtx.Run("maps a validation failure onto 400, not 500", func(subTestCtx *testing.T) {
		data := validTelemetry()
		data.DeviceID = ""

		recorder := postTelemetry(subTestCtx, server, data)

		assert.Equal(subTestCtx, http.StatusBadRequest, recorder.Code)
		assert.Equal(subTestCtx, "missing device ID", decodeError(subTestCtx, recorder))
	})

	// ---

	// NOTE: gin answers an unregistered *method* on a known path with 404 rather than 405, because
	// Engine.HandleMethodNotAllowed defaults to false. These pin the behaviour as currently
	// configured, not as HTTP would ideally have it.
	routingCases := []struct {
		Name   string
		Method string
		Path   string
	}{
		{Name: "rejects an unknown path", Method: http.MethodPost, Path: "/v1/nope"},
		{Name: "rejects an unversioned path", Method: http.MethodPost, Path: "/telemetry"},
		{
			Name:   "rejects a wrong method on a known path",
			Method: http.MethodGet,
			Path:   api.TelemetryPath,
		},
	}

	for _, testCase := range routingCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			recorder := doRequest(server, testCase.Method, testCase.Path, nil)

			assert.Equal(subTestCtx, http.StatusNotFound, recorder.Code)
		})
	}
}
