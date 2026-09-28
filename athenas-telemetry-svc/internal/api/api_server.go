// Package api provides the API server implementation for the telemetry service, including the
// gin.Engine setup and handler initialization.
package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/angiebrr/athenas-telemetry-svc/internal/telemetry"
)

// ================================================================================================

// ServerMode is used to control logging verbosity.
//
// We have our own enum here so callers outside of "api" don't need to know about gin / keep deps
// from leaking, and it also allows us to have stronger typing (since gin's modes are just strings)
type ServerMode string

const (
	// TestMode suppresses console output to keep test results clean
	TestMode ServerMode = gin.TestMode
	// DebugMode is the default mode and shows verbose debug logging + full stack traces
	DebugMode ServerMode = gin.DebugMode
	// ReleaseMode hides startup + debug logs
	ReleaseMode ServerMode = gin.ReleaseMode
)

// IsValid returns whether or not the current ServerMode is a valid one, since the underlying value
// is a string and really could be anything.
func (rMode ServerMode) IsValid() bool {
	switch rMode {
	case TestMode, DebugMode, ReleaseMode:
		return true
	default:
		return false
	}
}

// ------------------------------------------------------------------------------------------------

// OptionFn is the helper func for modifying the Options input struct
type OptionFn func(opts *Options)

// Options defines the args that can be passed into NewServer that can be modified by OptionFn
type Options struct {
	Mode ServerMode
}

// Options.Validate ensures that the input values are correct.
//
// Although we set defaults for optional params, we still validate everything since they
// could be overridden with input OptionFns
func (rOpts *Options) Validate() error {
	if !rOpts.Mode.IsValid() {
		return fmt.Errorf("invalid server mode [mode=%s]", rOpts.Mode)
	}

	return nil
}

// WithServerMode sets the mode for the server, which mostly just dictates what kinds of logs are
// output from the underlying gin engine.
//
// Note that updating this will update this globally for gin.
func WithServerMode(mode ServerMode) OptionFn {
	return func(opts *Options) {
		opts.Mode = mode
	}
}

// ------------------------------------------------------------------------------------------------

// Server encapsulates the main gin.Engine instance and any other configuration needed to host and
// serve the telemetry service
type Server struct {
	*gin.Engine

	// The internal domain logic service for managing telemetry
	telemetrySvc *telemetry.Service
}

// NewServer creates a Server instance by setting up the gin.Engine instance and initializing its
// handlers.
func NewServer(telemetrySvc *telemetry.Service, inOptFns ...OptionFn) (*Server, error) {
	opts := &Options{
		Mode: DebugMode,
	}
	for _, optFn := range inOptFns {
		optFn(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options to NewServer: %w", err)
	}

	// set the gin engine mode, which has to be done before creating an engine
	// (this is global)
	gin.SetMode(string(opts.Mode))

	// create the gin engine and set up handlers the engine will use for the server
	router := gin.Default()
	InitHandlers(router, telemetrySvc)

	return &Server{
		router,
		telemetrySvc,
	}, nil
}
