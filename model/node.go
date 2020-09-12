package model

// Node represents an URL like /blog that contains multiple
// pages, an overview page (ListPage) and child routes.
type Node struct {
	Children map[string]*Node
	Pages    []Page
	ListPage ListPage
}
