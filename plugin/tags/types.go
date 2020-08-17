package tags

import "github.com/verless/verless/model"

type indexPage struct {
	Meta *model.Meta
	Nav  *model.Nav
	*model.IndexPage
	Footer *model.Footer
}
