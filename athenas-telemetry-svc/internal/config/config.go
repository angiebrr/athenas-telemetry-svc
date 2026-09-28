// Package config holds the logic for deserializing env vars + config files into an AppConfig struct
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// ================================================================================================

const (
	// DefaultConfigPath is used if not overridden in config.InitEnv, which is a dotenv file in
	// the root directory
	DefaultConfigPath = ".env"

	// DefaultConfigType is the viper config type that's used if not overridden in config.InitEnv,
	// which is a dotenv file
	DefaultConfigType = "env"
)

// ------------------------------------------------------------------------------------------------

// OptionFn is the helper func for modifying the Options input struct
type OptionFn func(*Options)

// Options defines the args that can be passed into InitEnv that can be modified by OptionFn
type Options struct {
	ViperIns      *viper.Viper // The viper instance used for processing config + env vars
	AppConfigPath string       // The config path to pass into viper
	AppConfigType string       // The type of config to pass into viper
}

// Options.Validate ensures that the input values are correct.
//
// Although we set defaults for optional params, we still validate everything since they
// could be overridden with input OptionFns
func (rOpts *Options) Validate() error {
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

// WithViperIns sets the viper instance when processing the app's config + env vars.
// By default, a viper instance is created for you.
func WithViperIns(viperIns *viper.Viper) OptionFn {
	return func(opts *Options) {
		opts.ViperIns = viperIns
	}
}

// WithAppConfigPath sets the config path that is passed to viper to tell it where the config file
// should live.
//
// See DefaultConfigPath for what the default value is used if this is not called.
func WithAppConfigPath(cfgPath string) OptionFn {
	return func(opts *Options) {
		opts.AppConfigPath = cfgPath
	}
}

// WithAppConfigType sets the type of config that is passed to viper to tell it what type of config
// file should be expected.
//
// See DefaultConfigType for what the default value is used if this is not called.
func WithAppConfigType(cfgType string) OptionFn {
	return func(opts *Options) {
		opts.AppConfigType = cfgType
	}
}

// ------------------------------------------------------------------------------------------------

// InitEnv loads environment variables and dotenv files into a struct using viper.
//
// It accepts a slice of OptionFn funcs that can be used to override the default values.
func InitEnv(inOptFns ...OptionFn) (*AppConfig, error) {
	opts := &Options{
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
