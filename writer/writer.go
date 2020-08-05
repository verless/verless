package writer

import "github.com/verless/verless/model"

type Writer interface {
	Write(site model.Site)
}
