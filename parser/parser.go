package parser

import "github.com/verless/verless/model"

type Parser interface {
	ParsePage(src []byte) (model.Page, error)
}
