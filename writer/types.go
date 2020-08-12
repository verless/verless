package writer

import "github.com/verless/verless/model"

type page struct {
	Meta   *model.Meta
	Nav    *model.Nav
	Page   *model.Page
	Footer *model.Footer
}

type indexPage struct {
	Meta   *model.Meta
	Nav    *model.Nav
	Page   *model.IndexPage
	Footer *model.Footer
}
