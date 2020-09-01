package model

// Node represents an URL like /blog that contains multiple
// pages, an overview page (IndexPage) and child routes.
type Node struct {
	Children  map[string]*Node
	Pages     []Page
	IndexPage IndexPage
}
