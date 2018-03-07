package giraffe

import (
	"net/http"
)

var htmlContentType = []string{"text/html; charset=utf-8"}

// Render is actually a gin.render.Render
type Render struct {
	Giraffe *Giraffe
	Name    string
	Data    interface{}
}

// Render fulfills the gin.render.Render interface type
func (r Render) Render(w http.ResponseWriter) error {
	return r.Giraffe.render(w, r.Name, r.Data, true)
}

// WriteContentType fulfills the gin.render.Render interface type
func (r Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	contentType := header["Content-Type"]
	if len(contentType) == 0 {
		header["Content-Type"] = htmlContentType
	}
}
