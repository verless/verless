// Package theme provides functions for working with verless themes.
//
// Instead of juggling around with theme paths on your own, you should
// use the convenience functions from this package.
package theme

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/verless/verless/config"
)

const (
	TemplatesDir     string = "templates"
	CssDir           string = "css"
	JsDir            string = "js"
	Default          string = "default"
	PageTemplate     string = "page.html"
	ListPageTemplate string = "list-page.html"
	configFilename   string = "theme"
)

// Path returns the directory path for the theme with the given name
// inside the given path. Path does not ensure that the directory
// physically exists.
func Path(path, name string) string {
	return filepath.Join(path, config.ThemesDir, name)
}

// TemplatePath returns the template directory path of a given theme.
func TemplatePath(path, name string) string {
	return filepath.Join(Path(path, name), TemplatesDir)
}

// CssPath returns the css directory path of a given theme.
func CssPath(path, name string) string {
	return filepath.Join(Path(path, name), CssDir)
}

// JsPath returns the js directory path of a given theme.
func JsPath(path, name string) string {
	return filepath.Join(Path(path, name), JsDir)
}

// Exists determines whether a theme with the provided name inside
// the given path exists.
func Exists(path, name string) bool {
	if _, err := os.Stat(Path(path, name)); os.IsNotExist(err) {
		return false
	}
	return true
}

// Config represents a theme configuration. This is the configuration
// stored in the theme.yml file, which currently is not mandatory.
type Config struct {
	Version string
	Build   struct {
		Before []string
	}
}

// GetConfig returns the configuration stored in theme.yml of the theme
// with the given name inside the given path. Since theme.yml isn't
// mandatory, GetConfig returns an empty config if it doesn't exist.
func GetConfig(path, name string) (Config, error) {
	viper.AddConfigPath(Path(path, name))
	viper.SetConfigName(configFilename)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
		return cfg, nil
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// RunBeforeHooks executes all pre-build commands specified in the
// configuration of the theme with the specified name.
//
// Note that the command context directory is the the theme directory
// instead of the project directory.
func RunBeforeHooks(path, name string) error {
	cfg, err := GetConfig(path, name)
	if err != nil {
		return err
	}

	for _, beforeHook := range cfg.Build.Before {
		parts := strings.Split(beforeHook, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = Path(path, name)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
