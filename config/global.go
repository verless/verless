package config

var (
	// Version is injected when building a new release.
	Version string = "UNDEFINED"
)

const (
	// Filename (externally available as config.Filename) is the
	// name of the configuration file without extension.
	Filename string = "verless"

	// ContentDir is the directory for Markdown content.
	ContentDir string = "content"

	// TemplateDir is the directory for templates.
	TemplateDir string = "templates"

	// AssetDir is the directory for assets.
	AssetDir string = "assets"

	// OutputDir is the default output directory. This constant
	// is only used to set a default. The actual output directory
	// will be passed to the writer.
	OutputDir string = "target"

	// IndexFile is the filename used as directory index.
	IndexFile string = "index.html"

	// PageTpl is the filename for the template used for model.Page.
	// It has to live inside TemplateDir.
	PageTpl string = "page.html"

	// IndexPageTpl is the filename for the template used for
	// model.IndexPage. It has to live inside TemplateDir.
	IndexPageTpl string = "index-page.html"
)
