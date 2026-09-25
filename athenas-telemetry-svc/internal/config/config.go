package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// ================================================================================================

// FIXME: Add docstrings

const (
	DefaultConfigPath = ".env"
	DefaultConfigType = "env"
)

// ------------------------------------------------------------------------------------------------

type InitEnvOptFn func(*initEnvOptions)

type initEnvOptions struct {
	ViperIns      *viper.Viper
	AppConfigPath string
	AppConfigType string
}

func (rOpts *initEnvOptions) Validate() error {
	if rOpts.ViperIns == nil {
		return errors.New("viper instance is required")
	}
	if rOpts.AppConfigPath == "" {
		return errors.New("app config path is required")
	}
	if rOpts.AppConfigType == "" {
		return errors.New("app config type is required")
	}

	return nil
}

func WithViperIns(viperIns *viper.Viper) InitEnvOptFn {
	return func(opts *initEnvOptions) {
		opts.ViperIns = viperIns
	}
}

func WithAppConfigPath(cfgPath string) InitEnvOptFn {
	return func(opts *initEnvOptions) {
		opts.AppConfigPath = cfgPath
	}
}

func WithAppConfigType(cfgType string) InitEnvOptFn {
	return func(opts *initEnvOptions) {
		opts.AppConfigType = cfgType
	}
}

// ------------------------------------------------------------------------------------------------

// InitEnv loads environment variables and dotenv files into a struct using viper.
func InitEnv(inOptFns ...InitEnvOptFn) (*AppConfig, error) {
	opts := &initEnvOptions{
		ViperIns:      viper.New(),
		AppConfigPath: DefaultConfigPath,
		AppConfigType: DefaultConfigType,
	}
	for _, optFn := range inOptFns {
		optFn(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options to InitEnv: %w", err)
	}

	// set up app config defaults
	setAppConfigDefaults(opts.ViperIns)

	// config file type is dotenv, and is generally in the root directory
	opts.ViperIns.SetConfigType(opts.AppConfigType)
	opts.ViperIns.SetConfigFile(opts.AppConfigPath)

	// also load env vars into viper
	opts.ViperIns.AutomaticEnv()

	// after all that set up, actually read the values
	// note that we don't fail if the config file isn't found, since we may have used env vars instead
	if err := opts.ViperIns.ReadInConfig(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("unable to read config: %w", err)
		}
		fmt.Println("[warn] couldn't find config, will read from environment or use defaults")
	}

	// put the read values into a struct and make sure it's valid
	var parsedCfg AppConfig
	if err := opts.ViperIns.Unmarshal(&parsedCfg); err != nil {
		return nil, fmt.Errorf("unable to parse config with viper: %w", err)
	}
	if err := parsedCfg.Validate(); err != nil {
		return nil, err
	}

	return &parsedCfg, nil
}
