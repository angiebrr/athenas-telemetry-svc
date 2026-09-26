package config

import (
	"errors"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/spf13/viper"
)

// ================================================================================================

const (
	// EnvPort is the name of the viper config / env var value for the app's port
	EnvPort = "PORT"

	// EnvEnv is the name of the viper config / env var value for the app's environment
	EnvEnv = "ENV"
)

const (
	// DefaultPort is the default port value used for viper (see setAppConfigDefaults)
	DefaultPort = 8080
	// MinPortSystemThreshold is the min port value that can be used so system ports can't be used
	MinPortSystemThreshold = 1024
	// MaxPortThreshold is the largest port number that can be used in general
	MaxPortThreshold = 65535

	// DefaultEnv is the default env value used for viper (see setAppConfigDefaults)
	DefaultEnv = ValidEnvDev
	// ValidEnvDev is the name for the dev env
	ValidEnvDev = "dev"
	// ValidEnvProd is the name for the prod env
	ValidEnvProd = "prod"
)

// ------------------------------------------------------------------------------------------------

// DefaultAppConfig returns a config with the default values.
func DefaultAppConfig() AppConfig {
	return AppConfig{
		Env:  DefaultEnv,
		Port: DefaultPort,
	}
}

// setAppConfigDefaults sets defaults for all AppConfig values that have one using the given
// viper instance
func setAppConfigDefaults(viperIns *viper.Viper) {
	defaultCfg := DefaultAppConfig()
	viperIns.SetDefault(EnvPort, defaultCfg.Port)
	viperIns.SetDefault(EnvEnv, defaultCfg.Env)
}

// ------------------------------------------------------------------------------------------------

// AppConfig represents the possible values we can use for configuring the telemetry service via
// env vars and a dotenv file.
type AppConfig struct {
	Port int    `mapstructure:"PORT"` // The port thate the service is running on
	Env  string `mapstructure:"ENV"`  // The environment that the service is running in
}

// AppConfig.Validate returns whether or not the parsed config is valid
func (rCfg AppConfig) Validate() error {
	if rCfg.Port < MinPortSystemThreshold || rCfg.Port > MaxPortThreshold {
		errMsg := shared.FormatString(
			"{PORT_NAME} must be between {MIN} and {MAX} [{PORT_NAME}={PORT_VALUE}]",
			"PORT_NAME", EnvPort,
			"PORT_VALUE", rCfg.Port,
			"MIN", MinPortSystemThreshold,
			"MAX", MaxPortThreshold,
		)
		return errors.New(errMsg)
	}

	switch rCfg.Env {
	case ValidEnvDev, ValidEnvProd:
	default:
		errMsg := shared.FormatString(
			"{ENV_NAME} is must be either {DEV} {PROD} [{ENV_NAME}={ENV_VALUE}]",
			"ENV_NAME", EnvEnv,
			"ENV_VALUE", rCfg.Env,
			"DEV", ValidEnvDev,
			"PROD", ValidEnvProd,
		)
		return errors.New(errMsg)
	}

	return nil
}
