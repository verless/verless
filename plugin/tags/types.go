package tags

import "github.com/verless/verless/model"

// indexPage is an index page handed over to a template.
type indexPage struct {
	Meta *model.Meta
	Nav  *model.Nav
	*model.IndexPage
	Footer *model.Footer
}
