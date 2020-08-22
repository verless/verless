package parser

import (
	"time"

	"github.com/verless/verless/model"
)

const (
	dateFormat = "2006-01-02"
)

type (
	metadata map[string]interface{}
	assignFn func(val interface{})
)

func readMetadata(metadata metadata, page *model.Page) {
	readPrimitive(metadata["Title"], func(val interface{}) {
		page.Title = val.(string)
	})

	readPrimitive(metadata["Author"], func(val interface{}) {
		page.Author = val.(string)
	})

	readDate(metadata["Date"], func(val interface{}) {
		page.Date = val.(time.Time)
	})

	readList(metadata["Tags"], func(val interface{}) {
		page.Tags = append(page.Tags, val.(string))
	})

	readPrimitive(metadata["Img"], func(val interface{}) {
		page.Img = val.(string)
	})

	readPrimitive(metadata["Credit"], func(val interface{}) {
		page.Credit = val.(string)
	})

	readPrimitive(metadata["Description"], func(val interface{}) {
		page.Description = val.(string)
	})

	readList(metadata["Related"], func(val interface{}) {
		page.AddRelatedFQN(val.(string))
	})

	readPrimitive(metadata["Template"], func(val interface{}) {
		page.Template = val.(string)
	})

	readPrimitive(metadata["Hidden"], func(val interface{}) {
		page.SetHidden(val.(bool))
	})
}

func readPrimitive(field interface{}, assign assignFn) {
	if field == nil {
		return
	}

	assign(field)
}

func readDate(field interface{}, assign assignFn) {
	if field == nil {
		return
	}

	date, err := time.Parse(dateFormat, field.(string))
	if err != nil {
		panic(err)
	}

	assign(date)
}

func readList(field interface{}, assign assignFn) {
	if field == nil {
		return
	}

	list, _ := field.([]interface{})

	for i := 0; i < len(list); i++ {
		assign(list[i])
	}
}
