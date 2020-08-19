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
		Meta model.Meta
		Nav  struct {
			Items []struct {
				Label  string
				Target string
			}
			Overwrite bool
		}
		Footer struct {
			Items []struct {
				Label  string
				Target string
			}
			Overwrite bool
		}
	}
	Plugins []string
	Build   struct {
		Overwrite bool
	}
}

// HasPlugin checks if the configuration has enabled a given plugin.
func (c Config) HasPlugin(key string) bool {
	for _, plugin := range c.Plugins {
		if key == plugin {
			return true
		}
	}
	return false
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
