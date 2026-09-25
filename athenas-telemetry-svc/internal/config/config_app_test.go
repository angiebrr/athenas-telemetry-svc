package config_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/config"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

// FIXME: Add docstrings

func TestAppConfigValidate(testCtx *testing.T) {
	testCases := []struct {
		Name          string
		InputCfg      config.AppConfig
		ExpectedError string
	}{
		{
			Name:     "valid config: typical port",
			InputCfg: config.AppConfig{Port: 8080},
		},
		{
			Name:          "invalid config: system port",
			InputCfg:      config.AppConfig{Port: 1023},
			ExpectedError: config.EnvPort,
		},
		{
			Name:     "valid config: just outside system port",
			InputCfg: config.AppConfig{Port: 1024},
		},
		{
			Name:          "invalid config: outside max",
			InputCfg:      config.AppConfig{Port: 65546},
			ExpectedError: config.EnvPort, // Error should have the env var name in it
		},
		{
			Name:     "valid config: just outside system port",
			InputCfg: config.AppConfig{Port: 65535},
		},
	}

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			err := testCase.InputCfg.Validate()
			if testCase.ExpectedError != "" {
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				assert.NoError(subTestCtx, err)
			}
		})
	}
}
