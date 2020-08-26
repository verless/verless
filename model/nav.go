package model

// Nav represents the website's navigation.
type Nav struct {
	Items []NavItem
}

// NavItem represents an item in the navigation.
type NavItem struct {
	Label  string
	Target string
}
