package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/angiebrr/athenas-telemetry-svc/specifications"
)

// ================================================================================================

type IngestError struct {
	Err string `json:"error"`
}

func (rErr IngestError) Error() string {
	return rErr.Err
}

// ------------------------------------------------------------------------------------------------

type Driver struct {
	BaseURL *url.URL
	Client  *http.Client
}

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

func (rDriver *Driver) Ingest(data specifications.Telemetry) error {
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
