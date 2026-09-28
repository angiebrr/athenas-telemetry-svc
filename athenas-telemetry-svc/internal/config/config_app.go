package config

import (
	"errors"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/spf13/viper"
)

// ================================================================================================

const (
	// PortName is the name of the viper config / env var value for the app's port
	PortName = "PORT"

	// EnvName is the name of the viper config / env var value for the app's environment
	EnvName = "APP_ENV"
)

const (
	// DefaultPort is the default port value used for viper (see setAppConfigDefaults)
	DefaultPort = 8080
	// MinPortSystemThreshold is the min port value that can be used so system ports can't be used
	MinPortSystemThreshold = 1024
	// MaxPortThreshold is the largest port number that can be used in general
	MaxPortThreshold = 65535

	// DefaultEnv is the default env value used for viper (see setAppConfigDefaults)
	DefaultEnv = EnvDev
	// EnvDev is the name for the dev env
	EnvDev = "dev"
	// EnvProd is the name for the prod env
	EnvProd = "prod"
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
	viperIns.SetDefault(PortName, defaultCfg.Port)
	viperIns.SetDefault(EnvName, defaultCfg.Env)
}

// ------------------------------------------------------------------------------------------------

// AppConfig represents the possible values we can use for configuring the telemetry service via
// env vars and a dotenv file.
type AppConfig struct {
	Port int    `mapstructure:"PORT"`    // The port that the service is running on
	Env  string `mapstructure:"APP_ENV"` // The environment that the service is running in
}

// AppConfig.Validate returns whether or not the parsed config is valid
func (rCfg AppConfig) Validate() error {
	if rCfg.Port < MinPortSystemThreshold || rCfg.Port > MaxPortThreshold {
		errMsg := shared.FormatString(
			"{PORT_NAME} must be between {MIN} and {MAX}. [{PORT_NAME}={PORT_VALUE}]",
			"PORT_NAME", PortName,
			"PORT_VALUE", rCfg.Port,
			"MIN", MinPortSystemThreshold,
			"MAX", MaxPortThreshold,
		)
		return errors.New(errMsg)
	}

	switch rCfg.Env {
	case EnvDev, EnvProd:
	default:
		errMsg := shared.FormatString(
			"{ENV_NAME} must be either {DEV} or {PROD}. [{ENV_NAME}={ENV_VALUE}]",
			"ENV_NAME", EnvName,
			"ENV_VALUE", rCfg.Env,
			"DEV", EnvDev,
			"PROD", EnvProd,
		)
		return errors.New(errMsg)
	}

	return nil
}
