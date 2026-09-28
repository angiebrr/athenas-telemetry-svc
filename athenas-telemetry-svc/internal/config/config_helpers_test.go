package config_test

import "github.com/angiebrr/athenas-telemetry-svc/internal/config"

// ================================================================================================

// configWithPort copies the given cfg but with the given port value.
func configWithPort(cfg config.AppConfig, port int) config.AppConfig {
	cfg.Port = port
	return cfg
}

// configWithEnv copies the given cfg but with the given env value.
func configWithEnv(cfg config.AppConfig, env string) config.AppConfig {
	cfg.Env = env
	return cfg
}
