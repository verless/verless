package parser

import (
	"strings"
	"time"

	"github.com/verless/verless/model"
)

const (
	// dateFormat is the default date format expected for
	// the Date field in Markdown files.
	dateFormat = "2006-01-02"
)

type (
	// metadata represents a set of metadata.
	metadata map[string]interface{}

	// assignFn is a function that takes a value and assigns
	// that value to an enclosed struct field. For example:
	//
	//	setID := func(val interface{}) {
	//		user.ID = val.(int)
	//	}
	//
	//	func setter(setID assignFn) {
	//		setID(1)
	//	}
	//
	//	setter(setID)
	//
	// Type-asserting val to int in a closure avoids that the
	// setter function has to be re-written for each data type.
	assignFn func(val interface{})
)

// readMetadata reads values from a metadata map and assigns the
// values to the fields of a model.Page instance.
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
		name := val.(string)
		page.Tags = append(page.Tags, model.Tag{
			Name: name,
			Href: "/tags/" + strings.ToLower(name),
		})
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
		page.AddProvidedRelated(val.(string))
	})

	readPrimitive(metadata["Type"], func(val interface{}) {
		page.SetProvidedType(val.(string))
	})

	readPrimitive(metadata["Hidden"], func(val interface{}) {
		page.Hidden = val.(bool)
	})

	readPrimitive(metadata["Robots"], func(val interface{}) {
		page.Robots = val.(string)
	})
}

// readPrimitive converts a field to a primitive value and
// invokes the assignFn with that value.
func readPrimitive(field interface{}, assign assignFn) {
	if field == nil {
		return
	}

	assign(field)
}

// readPrimitive converts a field to a date and invokes the
// assignFn with that date.
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

// readList converts a field to a list and invokes the
// assignFn for each item in that list.
func readList(field interface{}, assign assignFn) {
	if field == nil {
		return
	}

	list, _ := field.([]interface{})

	for i := 0; i < len(list); i++ {
		assign(list[i])
	}
}
