package model

type Footer struct {
	Items []FooterItem
}

type FooterItem struct {
	Label  string
	Target string
}
