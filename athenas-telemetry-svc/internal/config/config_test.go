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

// writeTempDotEnv creates a dotenv file with the given contents and puts it in a temp directory
// that is cleaned up when all tests / subtests for testCtx are complete, and returns the path
// of the created file.
func writeTempDotEnv(testCtx *testing.T, contents string) string {
	testCtx.Helper()

	dir := testCtx.TempDir()
	configPath := filepath.Join(dir, ".env")

	// write the file using the typical permissions of dotenv files, which is only the owner of the
	// file can read and write it
	err := os.WriteFile(configPath, []byte(contents), 0600)
	require.NoError(testCtx, err)

	return configPath
}

// ------------------------------------------------------------------------------------------------

// TestConfigInit uses a mixture of setting env vars and creating temporary .env files to verify
// that parsing config values for our service works as expected.
//
// Most notably, we must verify the precedence of config parsing such that:
//
//	env vars > .env files > default values
func TestConfigInit(testCtx *testing.T) {
	validCfg := config.DefaultAppConfig()
	testCases := []struct {
		Name           string
		DotEnvContents string
		EnvVars        map[string]string
		ExpectedCfg    config.AppConfig
		ExpectedError  string
	}{
		{
			Name:        "no .env and no env vars, defaults used",
			ExpectedCfg: validCfg,
		},
		{
			Name:           "only .env sets PORT",
			DotEnvContents: "PORT=9090",
			ExpectedCfg:    configWithPort(validCfg, 9090),
		},
		{
			Name: "only env var sets PORT",
			EnvVars: map[string]string{
				config.EnvPort: "7070",
			},
			ExpectedCfg: configWithPort(validCfg, 7070),
		},
		{
			Name:           "env var takes precedence over .env",
			DotEnvContents: "PORT=9090",
			EnvVars: map[string]string{
				config.EnvPort: "7070",
			},
			ExpectedCfg: configWithPort(validCfg, 7070),
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
