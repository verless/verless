package model

// Route represents an URL like /blog that contains multiple
// pages, an overview page (IndexPage) and child routes.
type Route struct {
	Children  map[string]*Route
	Pages     []Page
	IndexPage IndexPage
}
