package packr

import (
	"github.com/gobuffalo/packr"
)

// Datasource uses a packr box for providing data
type Datasource struct {
	Box       packr.Box // The Box containing all the templates
	Extension string    // We need the extension to use convenience names
}

// LoadContent is needed to fulfill giraffe.Datasource interface type
func (d *Datasource) LoadContent(name string) (string, error) {
	return d.Box.MustString(name + d.Extension)
}
