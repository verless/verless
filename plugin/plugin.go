package plugin

import (
	"github.com/spf13/afero"
	"github.com/verless/verless/config"
	"github.com/verless/verless/model"
	"github.com/verless/verless/plugin/atom"
	"github.com/verless/verless/plugin/related"
	"github.com/verless/verless/plugin/tags"
)

// Plugin represents a built-in verless plugin.
type Plugin interface {
	// ProcessPage will be invoked after parsing the page. Must be safe
	// for concurrent usage.
	ProcessPage(page *model.Page) error
	// PreWrite will be invoked before writing the site.
	PreWrite(site *model.Site) error
	// PostWrite will be invoked after writing the site.
	PostWrite() error
}

// LoadAll returns a map of all available plugins. Each entry
// is a function that returns a fully initialized plugin instance.
func LoadAll(cfg *config.Config, fs afero.Fs, outputDir string) map[string]func() Plugin {
	return map[string]func() Plugin{
		"atom":    func() Plugin { return atom.New(&cfg.Site.Meta, fs, outputDir) },
		"related": func() Plugin { return related.New() },
		"tags":    func() Plugin { return tags.New() },
	}
}
