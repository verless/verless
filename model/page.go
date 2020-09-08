package model

import "time"

// Page represents a sub-page of the website.
type Page struct {
	Route       string
	ID          string
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
	Template    string

	providedRelated []string
	providedType    string
	hidden          bool
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

// Hidden describes if the page should be shown (false) or hidden (true).
func (p *Page) Hidden() bool {
	return p.hidden
}

// SetHidden shows (false) or hides (true) the page.
func (p *Page) SetHidden(hidden bool) {
	p.hidden = hidden
}

// IndexPage represents an overview page that is generated for
// each content sub-directory.
type IndexPage struct {
	Page
	Pages []*Page
}

// Type represents a page type.
type Type struct {
	Template string
}
