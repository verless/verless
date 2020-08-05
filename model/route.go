package model

type Route struct {
	Children  map[string]*Route
	Pages     []Page
	IndexPage IndexPage
}
