package config

import (
	"github.com/spf13/viper"
	"github.com/verless/verless/model"
)

type Config struct {
	Site struct {
		Meta model.Meta
		Nav  struct {
			Items []struct {
				Label  string
				Target string
			}
			Override bool
		}
		Footer struct {
			Items []struct {
				Label  string
				Target string
			}
			Override bool
		}
	}
}

func FromFile(path, filename string) (Config, error) {
	viper.AddConfigPath(path)
	// set the filename (without extension) this allows free usage of all formats viper accepts.
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
