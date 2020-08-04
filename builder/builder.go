package builder

import "github.com/verless/verless/model"

type Builder interface {
	RegisterPage(page model.Page) error
	Dispatch() (model.Site, error)
}
