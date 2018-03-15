package gin

import (
	"net/http"

	"github.com/benweidig/giraffe"

	"github.com/gin-gonic/gin/render"
)

// Giraffe is a thin wrapper around giraffe.Giraffe, so we cann add get around the
// non-local type no func adding thingie
type Giraffe struct {
	giraffe *giraffe.Giraffe
}

// New creates a new Giraffe but you have to bring your own config
func New(config giraffe.Config) *Giraffe {
	g := giraffe.New(config)
	return &Giraffe{g}
}

// Default creates a new Giraffe with sensible defaults
func Default() *Giraffe {
	g := giraffe.Default()
	return &Giraffe{g}
}

// Debug creates a new Giraffe for debugging
func Debug() *Giraffe {
	g := giraffe.Debug()
	return &Giraffe{g}
}

// Render is a helper struct that implements the gin.render.Render interface
type Render struct {
	Giraffe *giraffe.Giraffe
	Name    string
	Data    interface{}
}

// Render fulfills the gin.render.Render interface type
func (r Render) Render(w http.ResponseWriter) error {
	return r.Giraffe.Render(w, r.Name, r.Data, true)
}

var htmlContentType = []string{"text/html; charset=utf-8"}

// WriteContentType fulfills the gin.render.Render interface type
func (r Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	contentType := header["Content-Type"]
	if len(contentType) == 0 {
		header["Content-Type"] = htmlContentType
	}
}

// Instance fulfills the gin.render.HTMLRender interface type.
// Giraffe isn't specifically for gin, but we want to fullfil the interface
func (g *Giraffe) Instance(name string, data interface{}) render.Render {
	return Render{
		Giraffe: g.giraffe,
		Name:    name,
		Data:    data,
	}
}
