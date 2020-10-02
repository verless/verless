package model

import "time"

// Page represents a sub-page of the website.
type Page struct {
	Route       string
	ID          string
	Href        string
	Title       string
	Author      string
	Date        time.Time
	Tags        []string
	Img         string
	Credit      string
	Description string
	Content     string
	Related     []*Page
	Type        *Type
	Hidden      bool

	providedRelated []string
	providedType    string
}

// ProvidedRelated returns all Fully Qualified Name URIs related to the page.
func (p *Page) ProvidedRelated() []string {
	return p.providedRelated
}

// AddProvidedRelated adds a new Fully Qualified Name URI to the page.
func (p *Page) AddProvidedRelated(relatedFQN string) {
	p.providedRelated = append(p.providedRelated, relatedFQN)
}

// ProvidedType returns the user-provided page type.
func (p *Page) ProvidedType() string {
	return p.providedType
}

// SetProvidedType sets the user-provided page type.
func (p *Page) SetProvidedType(providedType string) {
	p.providedType = providedType
}

// ListPage represents an overview page that is generated for
// each content sub-directory.
type ListPage struct {
	Page
	Pages []*Page
	Route string
}

// Type represents a page type.
type Type struct {
	Template string
}
