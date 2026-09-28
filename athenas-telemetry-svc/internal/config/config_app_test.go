package config_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/config"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

// TestAppConfigValidate verifies that the validation helper for AppConfig accepts and rejects
// values as expected
func TestAppConfigValidate(testCtx *testing.T) {
	validCfg := config.DefaultAppConfig()
	testCases := []struct {
		Name          string
		InputCfg      config.AppConfig
		ExpectedError string
	}{
		{
			Name:     "valid config",
			InputCfg: validCfg,
		},
		{
			Name:          "invalid config: under min port",
			InputCfg:      configWithPort(validCfg, config.MinPortSystemThreshold-1),
			ExpectedError: config.PortName,
		},
		{
			Name:     "valid config: at min port",
			InputCfg: configWithPort(validCfg, config.MinPortSystemThreshold),
		},
		{
			Name:          "invalid config: above max port",
			InputCfg:      configWithPort(validCfg, config.MaxPortThreshold+1),
			ExpectedError: config.PortName,
		},
		{
			Name:     "valid config: at max port",
			InputCfg: configWithPort(validCfg, config.MaxPortThreshold),
		},
		{
			Name:          "invalid config: bad environment",
			InputCfg:      configWithEnv(validCfg, "venus"),
			ExpectedError: config.EnvName,
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
