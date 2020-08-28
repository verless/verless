package model

// Footer represents the website's footer.
type Footer struct {
	Items []FooterItem
}

// FooterItem represents an item in the footer.
type FooterItem struct {
	Label  string
	Target string
}
