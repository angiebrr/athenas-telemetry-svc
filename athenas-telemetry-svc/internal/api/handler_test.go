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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

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

// ================================================================================================

type HandlerTestSuite struct {
	suite.Suite
	server *api.Server
}

// HandlerTestSuite.SetupSuite builds the one server instance shared by every test in the suite.
//
// NOTE: The server holds no request state as of M1, so sharing it is safe. Once it owns a dispatch
// ring this has to move to SetupTest so cases can't leak state into each other.
func (rSuite *HandlerTestSuite) SetupSuite() {
	// quiet gin debug logs during testing
	gin.SetMode(gin.TestMode)

	// create a single server instance for all tests in this suite
	rSuite.server = api.NewServer()
	require.NotNil(rSuite.T(), rSuite.server, "server should not be nil")
}

// ------------------------------------------------------------------------------------------------

// HandlerTestSuite.doRequest sends one request directly into the gin engine and returns the recorded
// response.
func (rSuite *HandlerTestSuite) doRequest(
	method string,
	path string,
	body io.Reader,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	rSuite.server.ServeHTTP(recorder, request)

	return recorder
}

// HandlerTestSuite.postTelemetry serializes the given telemetry and POSTs it to the ingest endpoint.
func (rSuite *HandlerTestSuite) postTelemetry(data models.Telemetry) *httptest.ResponseRecorder {
	payload, err := json.Marshal(data)
	require.NoError(rSuite.T(), err, "test telemetry should serialize")

	return rSuite.doRequest(http.MethodPost, api.TelemetryPath, bytes.NewReader(payload))
}

// HandlerTestSuite.decodeError pulls the message out of the handler's JSON error envelope.
func (rSuite *HandlerTestSuite) decodeError(recorder *httptest.ResponseRecorder) string {
	var response errorResponse

	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(rSuite.T(), err, "error responses should carry a JSON error envelope")

	return response.Error
}

// ================================================================================================

// HandlerTestSuite.TestMalformedPayload verifies that a body which isn't valid JSON is rejected at
// the binding step, before the domain is ever reached.
//
// This is the only test in the repo that exercises the ShouldBindJSON failure branch — the spec
// can't reach it, since its interface speaks models.Telemetry rather than bytes.
func (rSuite *HandlerTestSuite) TestMalformedPayload() {
	recorder := rSuite.doRequest(
		http.MethodPost,
		api.TelemetryPath,
		bytes.NewBufferString(`{"device_id":`),
	)

	rSuite.Equal(http.StatusBadRequest, recorder.Code)
	rSuite.NotEmpty(rSuite.decodeError(recorder), "a malformed body should still report an error")
}

// HandlerTestSuite.TestValidPayload verifies the accepted path answers 202 Accepted rather than a
// bare 200, which is the distinction the spec cannot make on its own.
func (rSuite *HandlerTestSuite) TestValidPayload() {
	data := validTelemetry()
	recorder := rSuite.postTelemetry(data)

	rSuite.Equal(http.StatusAccepted, recorder.Code)
	rSuite.Empty(recorder.Body.String(), "an accepted ingest should not return a body")
}

// HandlerTestSuite.TestInvalidPayload verifies a domain validation failure maps onto 400 rather than
// the 500 the handler falls back to for internal errors.
//
// This is the concrete, HTTP-specific half of the error classification that the spec asserts only
// loosely (as "some error") while Ingest still has a single failure mode.
func (rSuite *HandlerTestSuite) TestInvalidPayload() {
	data := validTelemetry()
	data.DeviceID = ""

	recorder := rSuite.postTelemetry(data)

	rSuite.Equal(http.StatusBadRequest, recorder.Code)
	rSuite.Equal("missing device ID", rSuite.decodeError(recorder))
}

// HandlerTestSuite.TestInvalidPathOrMethod verifies that only POST /v1/telemetry is routed.
//
// NOTE: gin answers an unregistered *method* on a known path with 404 rather than 405, because
// Engine.HandleMethodNotAllowed defaults to false. This pins the behaviour as currently configured,
// not as HTTP would ideally have it.
func (rSuite *HandlerTestSuite) TestInvalidPathOrMethod() {
	cases := []struct {
		Name   string
		Method string
		Path   string
	}{
		{Name: "unknown path", Method: http.MethodPost, Path: "/v1/nope"},
		{Name: "unversioned path", Method: http.MethodPost, Path: "/telemetry"},
		{Name: "wrong method on a known path", Method: http.MethodGet, Path: api.TelemetryPath},
	}

	for _, testCase := range cases {
		rSuite.Run(testCase.Name, func() {
			recorder := rSuite.doRequest(testCase.Method, testCase.Path, nil)

			rSuite.Equal(http.StatusNotFound, recorder.Code)
		})
	}
}

// ================================================================================================

func TestHandlerTestSuite(testCtx *testing.T) {
	suite.Run(testCtx, new(HandlerTestSuite))
}
