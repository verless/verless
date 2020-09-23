package config

var (
	// GitTag is injected when building a new release.
	GitTag string = "UNDEFINED"
	// GitCommit stores the latest Git commit.
	GitCommit string = "UNKNOWN"
)

const (
	// Filename is the name of the config file without extension.
	Filename string = "verless"

	// ContentDir is the directory for Markdown content.
	ContentDir string = "content"

	// ThemesDir is the directory for verless themes.
	ThemesDir string = "themes"

	// TemplateDir is the directory for templates inside ThemesDir.
	TemplateDir string = "templates"

	// CSSDir is the directory for CSS files.
	CSSDir string = "css"

	// JSDir is the directory for JavaScript files.
	JSDir string = "js"

	// DefaultTheme is the name of the default theme.
	DefaultTheme string = "default"

	// CSSFile is the default CSS file expected in CSSDir.
	CSSFile string = "style.css"

	// StaticDir is the directory for static files.
	StaticDir string = "static"

	// OutputDir is the default output directory.
	OutputDir string = "target"

	// IndexFile is the filename used as directory index.
	IndexFile string = "index.html"

	// PageTpl is the template file used for model.Page.
	PageTpl string = "page.html"

	// ListPageTpl is the template file used for model.ListPage.
	ListPageTpl string = "list-page.html"

	// ListPageID is the ID for custom list pages that overwrite
	// a auto-generated list page.
	ListPageID string = "index"
)
