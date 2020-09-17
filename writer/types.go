package writer

import "github.com/verless/verless/model"

// page is a wrapper for Page-related templates. It gets directly
// passed to the template and allows the navigation to be accessed
// with `.Nav`, for example.
type page struct {
	Meta   *model.Meta
	Nav    *model.Nav
	Page   *model.Page
	Footer *model.Footer
}

// listPage is a wrapper for ListPage-related templates.
type listPage struct {
	Meta *model.Meta
	Nav  *model.Nav
	*model.ListPage
	Footer *model.Footer
}
