package theme

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"

	"github.com/verless/verless/config"
)

const configFilename string = "theme"

func Path(path, name string) string {
	return filepath.Join(path, config.ThemesDir, name)
}

func TemplateDir(path, name string) string {
	return filepath.Join(Path(path, name), config.TemplateDir)
}

func CssDir(path, name string) string {
	return filepath.Join(Path(path, name), config.CssDir)
}

func JsDir(path, name string) string {
	return filepath.Join(Path(path, name), config.JsDir)
}

func Exists(path, name string) bool {
	if _, err := os.Stat(Path(path, name)); os.IsNotExist(err) {
		return false
	}
	return true
}

type Config struct {
	Version string
	Build   struct {
		Before []string
	}
}

func GetConfig(path, name string) (Config, error) {
	viper.AddConfigPath(Path(path, name))
	viper.SetConfigName(configFilename)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
