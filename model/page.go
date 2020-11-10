package model

import "time"

const (
	customListPageID string = "index"
)

type Tag struct {
	Name string
	Href string
}

func (t Tag) String() string {
	return t.Name
}

// Page represents a sub-page of the website.
type Page struct {
	Route       string
	ID          string
	Href        string
	Title       string
	Author      string
	Date        time.Time
	Tags        []Tag
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

// IsCustomListPage returns whether the page is a custom list page that has
// been created from a file called index.md in a content directory.
func (p *Page) IsCustomListPage() bool {
	return p.ID == customListPageID
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
}

// Type represents a page type.
type Type struct {
	Template string
}
