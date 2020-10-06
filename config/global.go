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

	// GeneratedDir is the directory which can be used by hook-commands
	// and which gets ignored by the serve command.
	// The directory can exist in each theme directory and in the StaticDir.
	GeneratedDir string = "generated"

	// StaticDir is the directory for static files.
	StaticDir string = "static"

	// OutputDir is the default output directory.
	OutputDir string = "target"
)
