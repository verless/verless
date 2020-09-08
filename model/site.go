package model

// walkFn is invoked by WalkTree for each node in the route tree.
type walkFn func(node *Node) error

// Site represents the actual website. The site model is generated
// and populated with data and content during the website build.
//
// Any build.Writer implementation is capable of rendering this
// model as a static website.
type Site struct {
	Meta   Meta
	Nav    Nav
	Root   *Node
	Footer Footer
}
