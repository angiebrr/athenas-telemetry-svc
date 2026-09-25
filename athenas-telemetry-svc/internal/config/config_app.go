package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ================================================================================================

// FIXME: Add docstrings

const (
	EnvPort = "PORT"
)

const (
	DefaultPort            = 8080
	MinPortSystemThreshold = 1024
	MaxPortThreshold       = 65535
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
		return fmt.Errorf("%s must be between 1024 and 65535 [%s=%d]", EnvPort, EnvPort, rCfg.Port)
	}

	return nil
}

// ------------------------------------------------------------------------------------------------

func setAppConfigDefaults(viperIns *viper.Viper) {
	viperIns.SetDefault(EnvPort, DefaultPort)
}
