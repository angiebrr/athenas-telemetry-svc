package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/angiebrr/athenas-telemetry-svc/models"
)

// ================================================================================================

// IngestError is a helper used to deserialize errors reported by the telemetry HTTP service
type IngestError struct {
	Err string `json:"error"`
}

// IngestError.Error simply returns the deserialized error message that was in the error key
func (rErr IngestError) Error() string {
	return rErr.Err
}

// ------------------------------------------------------------------------------------------------

// Driver is used by acceptance tests to make calls to our telemetry HTTP server
type Driver struct {
	BaseURL *url.URL
	Client  *http.Client
}

// NewDriver sets up the base URL and HTTP client and creates a new Driver instance used for
// acceptance testing
func NewDriver(baseURL string, client *http.Client) (*Driver, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to make http driver: %w", err)
	}

	return &Driver{
		BaseURL: url,
		Client:  client,
	}, nil
}

// Driver.Ingest calls the POST /v1/telemetry HTTP endpoint with the given telemetry data, and
// reports the error if there is any.
func (rDriver *Driver) Ingest(data models.Telemetry) error {
	// set up the proper endpoint, starting with the base URL
	ingestURL := *rDriver.BaseURL
	ingestURL.Path = "v1/telemetry"

	// serialize our telemetry payload
	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("unable to serialize telemetry data for Ingest(): %w", err)
	}
	ingestBody := bytes.NewBuffer(jsonBytes)

	// call our POST telemetry endpoint with the payload
	resp, err := rDriver.Client.Post(ingestURL.String(), "application/json", ingestBody)
	if err != nil {
		return fmt.Errorf("failed to POST at telemetry endpoint: %w", err)
	}
	defer resp.Body.Close()

	// if not successful, try to get an error from it
	if resp.StatusCode != http.StatusAccepted {
		var ingestErr IngestError
		if err := json.NewDecoder(resp.Body).Decode(&ingestErr); err != nil {
			return fmt.Errorf("failed to read response body from POST telemetry endpoint: %w", err)
		}
		return ingestErr
	}

	return nil
}
