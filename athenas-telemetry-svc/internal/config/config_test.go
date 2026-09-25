package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ================================================================================================

// FIXME: Add docstrings

func writeTempDotEnv(testCtx *testing.T, contents string) string {
	testCtx.Helper()

	dir := testCtx.TempDir()
	configPath := filepath.Join(dir, ".env")

	err := os.WriteFile(configPath, []byte(contents), 0600)
	require.NoError(testCtx, err)

	return configPath
}

// ------------------------------------------------------------------------------------------------

func TestConfigInit(testCtx *testing.T) {
	testCases := []struct {
		Name           string
		DotEnvContents string
		EnvVars        map[string]string
		ExpectedCfg    config.AppConfig
		ExpectedError  string
	}{
		{
			Name:        "no .env and no env vars, defaults used",
			ExpectedCfg: config.AppConfig{Port: config.DefaultPort},
		},
		{
			Name:           "only .env sets PORT",
			DotEnvContents: "PORT=9090",
			ExpectedCfg:    config.AppConfig{Port: 9090},
		},
		{
			Name: "only env var sets PORT",
			EnvVars: map[string]string{
				config.EnvPort: "7070",
			},
			ExpectedCfg: config.AppConfig{Port: 7070},
		},
		{
			Name:           "env var takes precedence over .env",
			DotEnvContents: "PORT=9090",
			EnvVars: map[string]string{
				config.EnvPort: "7070",
			},
			ExpectedCfg: config.AppConfig{Port: 7070},
		},
		{
			Name: "env vars sets invalid port",
			EnvVars: map[string]string{
				config.EnvPort: "100",
			},
			ExpectedError: config.EnvPort, // expect an error that has the port env var name in it
		},
	}

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			cfgPath := config.DefaultConfigPath

			// setup .env files or env var values based on what the test case configured
			if testCase.DotEnvContents != "" {
				cfgPath = writeTempDotEnv(subTestCtx, testCase.DotEnvContents)
			}
			for envName, envValue := range testCase.EnvVars {
				subTestCtx.Setenv(envName, envValue)
			}

			// call the func we're testing, InitEnv
			cfg, err := config.InitEnv(
				config.WithAppConfigPath(cfgPath),
			)

			// verify that it either errors like we're expected or has the parsed config that we're expecting
			if testCase.ExpectedError != "" {
				assert.ErrorContains(subTestCtx, err, testCase.ExpectedError)
			} else {
				require.NoError(subTestCtx, err)
				require.NotNil(subTestCtx, cfg)

				got := *cfg
				want := testCase.ExpectedCfg
				assert.Equal(subTestCtx, want, got)
			}
		})
	}
}
