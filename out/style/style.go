// Package style provides icons, emojis and other stylized items
// that can be printed to the standard output device.
package style

// Emoji represents an emoji code point.
type Emoji string

// Enumerations for all supported emojis. This package uses the
// GitHub and Slack shortcodes as enumeration names.
const (
	None Emoji = ""
	Tada       = "🎉"
)
