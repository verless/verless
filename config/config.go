// Package config provides configuration-related types and functions.
package config

import (
	"github.com/spf13/viper"
	"github.com/verless/verless/model"
)

// Config represents the user configuration stored in verless.yml.
type Config struct {
	Version string
	Site    struct {
		Meta   model.Meta
		Nav    model.Nav
		Footer model.Footer
	}
	Plugins []string
	Theme   string
	Types   map[string]*model.Type
	Build   struct {
		Overwrite bool
	}
}

// FromFile looks for a configuration file and converts it to a Config.
func FromFile(path, filename string) (Config, error) {
	viper.AddConfigPath(path)
	// Set the filename without extension to allow all supported formats.
	viper.SetConfigName(filename)

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
