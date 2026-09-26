package config

import (
	"errors"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/spf13/viper"
)

// ================================================================================================

// FIXME: Add docstrings

const (
	// EnvPort is the name of the viper config / env var value for PORT
	EnvPort = "PORT"
)

const (
	// DefaultPort is the default port value used for viper (see setAppConfigDefaults)
	DefaultPort = 8080
	// MinPortSystemThreshold is the min port value that can be used so system ports can't be used
	MinPortSystemThreshold = 1024
	// MaxPortThreshold is the largest port number that can be used in general
	MaxPortThreshold = 65535
)

// ------------------------------------------------------------------------------------------------

// AppConfig represents the possible values we can use for configuring the telemetry service via
// env vars and a dotenv file.
type AppConfig struct {
	Port int `mapstructure:"PORT"`
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

	return nil
}

// ------------------------------------------------------------------------------------------------

// setAppConfigDefaults sets defaults for all AppConfig values that have one using the given
// viper instance
func setAppConfigDefaults(viperIns *viper.Viper) {
	viperIns.SetDefault(EnvPort, DefaultPort)
}
