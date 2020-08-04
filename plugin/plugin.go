package plugin

import "github.com/verless/verless/model"

type Plugin interface {
	ProcessPage(page *model.Page) error
	Finalize() error
}
