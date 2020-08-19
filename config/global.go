package config

var (
	// Version is injected when building a new release.
	Version string = "UNDEFINED"
)

const (
	// Filename is the name of the config file without extension.
	Filename string = "verless"

	// ContentDir is the directory for Markdown content.
	ContentDir string = "content"

	// TemplateDir is the directory for templates.
	TemplateDir string = "templates"

	// AssetDir is the directory for assets.
	AssetDir string = "assets"

	// OutputDir is the default output directory.
	OutputDir string = "target"

	// IndexFile is the filename used as directory index.
	IndexFile string = "index.html"

	// PageTpl is the template file used for model.Page.
	PageTpl string = "page.html"

	// IndexPageTpl is the template file used for model.IndexPage.
	IndexPageTpl string = "index-page.html"
)
